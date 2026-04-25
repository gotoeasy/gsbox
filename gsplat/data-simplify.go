package gsplat

import (
	"gsbox/cmn"
	"math"
	"runtime"
	"sort"
	"sync"
)

// 参数
const (
	gridSizeFactor     float32 = 2.5 // 网格大小系数，通常2到3，越大简化程度越激进速度越慢
	baseMergeThreshold float32 = -3.0
	scaleMinRatio      float32 = 0.5
	scaleMaxRatio      float32 = 3.0
	blockSize          float64 = 128.0
	invBlock           float64 = 1.0 / blockSize
)

type gridKey struct{ x, y, z int32 }

type cellInfo struct {
	gx, gy, gz int32
	idx        int32
}

type flatHashEntry struct {
	key  gridKey
	vals []int32
}

type flatHash struct {
	entries []flatHashEntry
	mask    int
}

// 高斯简化
func Simplify(splats []*SplatData) []*SplatData {
	n := len(splats)
	if n == 0 {
		return splats
	}

	var validSplats []*SplatData
	var avgScale float32
	for _, s := range splats {
		if s.ColorA < 20 {
			continue
		}
		alpha := fastOpacity(s.ColorA)
		s.Temp2Float32 = float32(math.Exp(float64(s.ScaleX)))
		s.Temp3Float32 = float32(math.Exp(float64(s.ScaleY)))
		s.Temp4Float32 = float32(math.Exp(float64(s.ScaleZ)))
		s.Temp1Float32 = s.Temp2Float32 * s.Temp3Float32 * s.Temp4Float32 * alpha

		validSplats = append(validSplats, s)
		avgScale += (s.Temp2Float32 + s.Temp3Float32 + s.Temp4Float32) / 3.0
	}
	if len(validSplats) == 0 {
		return validSplats
	}

	avgScale /= float32(len(validSplats))
	splats = validSplats

	sort.Slice(splats, func(i, j int) bool {
		return splats[i].Temp1Float32 > splats[j].Temp1Float32
	})

	type blockKey struct{ x, y, z int32 }
	type blockEntry struct{ list []*SplatData }
	blockMap := make(map[blockKey]*blockEntry)
	for _, s := range splats {
		bk := blockKey{
			int32(math.Floor(float64(s.PositionX) * invBlock)),
			int32(math.Floor(float64(s.PositionY) * invBlock)),
			int32(math.Floor(float64(s.PositionZ) * invBlock)),
		}
		if blockMap[bk] == nil {
			blockMap[bk] = &blockEntry{list: make([]*SplatData, 0, 64)}
		}
		blockMap[bk].list = append(blockMap[bk].list, s)
	}

	workers := runtime.NumCPU()
	taskCh := make(chan *blockEntry, len(blockMap))
	resCh := make(chan []*SplatData, len(blockMap))
	var wg sync.WaitGroup

	gridSize := avgScale * gridSizeFactor
	if gridSize < 1e-3 {
		gridSize = 0.01
	}

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for b := range taskCh {
				resCh <- processChunk(b.list, gridSize)
			}
		}()
	}

	go func() {
		for _, b := range blockMap {
			taskCh <- b
		}
		close(taskCh)
		wg.Wait()
		close(resCh)
	}()

	final := make([]*SplatData, 0, n>>1)
	for r := range resCh {
		final = append(final, r...)
	}
	return final
}

func newFlatHash(cap int) *flatHash {
	if cap < 8 {
		cap = 8
	}
	size := 1
	for size < int(float32(cap)*1.5) {
		size <<= 1
	}
	return &flatHash{
		entries: make([]flatHashEntry, size),
		mask:    size - 1,
	}
}

func (h *flatHash) add(k gridKey, v int32) {
	hash := int32(k.x)*73856093 + int32(k.y)*19349663 + int32(k.z)*83492791
	idx := int(hash) & h.mask
	for range 8 {
		ent := &h.entries[idx]
		if ent.vals == nil {
			ent.key = k
			ent.vals = append(ent.vals, v)
			return
		}
		if ent.key == k {
			ent.vals = append(ent.vals, v)
			return
		}
		idx = (idx + 1) & h.mask
	}
}

func (h *flatHash) get(k gridKey) []int32 {
	hash := int32(k.x)*73856093 + int32(k.y)*19349663 + int32(k.z)*83492791
	idx := int(hash) & h.mask
	for range 8 {
		ent := &h.entries[idx]
		if ent.key == k {
			return ent.vals
		}
		if ent.vals == nil {
			return nil
		}
		idx = (idx + 1) & h.mask
	}
	return nil
}

var splatPool = sync.Pool{
	New: func() any { return &SplatData{} },
}

func allocSplat() *SplatData {
	return splatPool.Get().(*SplatData)
}

func abs32(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

func fastOpacity(a uint8) float32 {
	if a <= 2 {
		return 0.001
	}
	if a >= 253 {
		return 0.999
	}
	return cmn.EncodeSplatOpacityF32(a)
}

func similarity(a, b *SplatData) float32 {
	dx := a.PositionX - b.PositionX
	dy := a.PositionY - b.PositionY
	dz := a.PositionZ - b.PositionZ
	distSq := dx*dx + dy*dy + dz*dz

	dr := float32(a.ColorR-b.ColorR) / 255.0
	dg := float32(a.ColorG-b.ColorG) / 255.0
	db := float32(a.ColorB-b.ColorB) / 255.0
	colorSq := dr*dr + dg*dg + db*db

	scaleDiff := abs32(a.Temp2Float32 - b.Temp2Float32)
	scaleSum := a.Temp2Float32 + b.Temp2Float32 + 1e-6

	raw := -(distSq + colorSq*0.5 + (scaleDiff/scaleSum)*2.0)
	return raw
}

func outer(x, y, z float32) [3][3]float32 {
	return [3][3]float32{
		{x * x, x * y, x * z},
		{y * x, y * y, y * z},
		{z * x, z * y, z * z},
	}
}

func matAdd(a, b [3][3]float32) [3][3]float32 {
	return [3][3]float32{
		{a[0][0] + b[0][0], a[0][1] + b[0][1], a[0][2] + b[0][2]},
		{a[1][0] + b[1][0], a[1][1] + b[1][1], a[1][2] + b[1][2]},
		{a[2][0] + b[2][0], a[2][1] + b[2][1], a[2][2] + b[2][2]},
	}
}

func matScale(a [3][3]float32, s float32) [3][3]float32 {
	return [3][3]float32{
		{a[0][0] * s, a[0][1] * s, a[0][2] * s},
		{a[1][0] * s, a[1][1] * s, a[1][2] * s},
		{a[2][0] * s, a[2][1] * s, a[2][2] * s},
	}
}

func quatToMat3(w, x, y, z float32) [3][3]float32 {
	xx, yy, zz := x*x, y*y, z*z
	return [3][3]float32{
		{1 - 2*(yy+zz), 2 * (x*y - z*w), 2 * (x*z + y*w)},
		{2 * (x*y + z*w), 1 - 2*(xx+zz), 2 * (y*z - x*w)},
		{2 * (x*z - y*w), 2 * (y*z + x*w), 1 - 2*(xx+yy)},
	}
}

func buildSigma(s *SplatData) [3][3]float32 {
	qw, qx, qy, qz := cmn.NormalizeRotationsUint8F32(s.RotationW, s.RotationX, s.RotationY, s.RotationZ)
	R := quatToMat3(qw, qx, qy, qz)
	D := [3][3]float32{
		{s.Temp2Float32 * s.Temp2Float32, 0, 0},
		{0, s.Temp3Float32 * s.Temp3Float32, 0},
		{0, 0, s.Temp4Float32 * s.Temp4Float32},
	}

	var res [3][3]float32
	for i := range 3 {
		for j := 0; j < 3; j++ {
			res[i][j] = R[i][0]*D[0][0]*R[j][0] + R[i][1]*D[1][1]*R[j][1] + R[i][2]*D[2][2]*R[j][2]
		}
	}
	return res
}

func eigenDecomp(A [3][3]float32) ([3]float32, [3][3]float32) {
	V := [3][3]float32{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
	for range 10 {
		p, q := 0, 1
		maxVal := float32(math.Abs(float64(A[0][1])))
		if math.Abs(float64(A[0][2])) > float64(maxVal) {
			p, q = 0, 2
			maxVal = float32(math.Abs(float64(A[0][2])))
		}
		if math.Abs(float64(A[1][2])) > float64(maxVal) {
			p, q = 1, 2
		}
		if math.Abs(float64(A[p][q])) < 1e-6 {
			break
		}

		theta := 0.5 * float32(math.Atan2(float64(2*A[p][q]), float64(A[q][q]-A[p][p])))
		c, s := float32(math.Cos(float64(theta))), float32(math.Sin(float64(theta)))

		for i := range 3 {
			api, aqi := A[p][i], A[q][i]
			A[p][i], A[q][i] = c*api-s*aqi, s*api+c*aqi
		}
		for i := range 3 {
			aip, aiq := A[i][p], A[i][q]
			A[i][p], A[i][q] = c*aip-s*aiq, s*aip+c*aiq
		}
		for i := range 3 {
			vip, viq := V[i][p], V[i][q]
			V[i][p], V[i][q] = c*vip-s*viq, s*vip+c*viq
		}
	}

	return [3]float32{A[0][0], A[1][1], A[2][2]}, V
}

func matToQuat(m [3][3]float32) (float32, float32, float32, float32) {
	trace := m[0][0] + m[1][1] + m[2][2]
	if trace > 0 {
		s := float32(math.Sqrt(float64(trace+1.0))) * 2
		return 0.25 * s, (m[2][1] - m[1][2]) / s, (m[0][2] - m[2][0]) / s, (m[1][0] - m[0][1]) / s
	}
	return 1, 0, 0, 0
}

func mergeTwo(a, b *SplatData) *SplatData {
	total := a.Temp1Float32 + b.Temp1Float32
	if total < 1e-6 {
		return a
	}
	invTotal := 1.0 / total
	o := allocSplat()

	o.PositionX = (a.PositionX*a.Temp1Float32 + b.PositionX*b.Temp1Float32) * invTotal
	o.PositionY = (a.PositionY*a.Temp1Float32 + b.PositionY*b.Temp1Float32) * invTotal
	o.PositionZ = (a.PositionZ*a.Temp1Float32 + b.PositionZ*b.Temp1Float32) * invTotal

	o.ColorR = uint8((float32(a.ColorR)*a.Temp1Float32 + float32(b.ColorR)*b.Temp1Float32) * invTotal)
	o.ColorG = uint8((float32(a.ColorG)*a.Temp1Float32 + float32(b.ColorG)*b.Temp1Float32) * invTotal)
	o.ColorB = uint8((float32(a.ColorB)*a.Temp1Float32 + float32(b.ColorB)*b.Temp1Float32) * invTotal)
	o.ColorA = uint8((float32(a.ColorA)*a.Temp1Float32 + float32(b.ColorA)*b.Temp1Float32) * invTotal)

	sa, sb := buildSigma(a), buildSigma(b)
	dax, day, daz := a.PositionX-o.PositionX, a.PositionY-o.PositionY, a.PositionZ-o.PositionZ
	dbx, dby, dbz := b.PositionX-o.PositionX, b.PositionY-o.PositionY, b.PositionZ-o.PositionZ

	sigma := matAdd(matScale(sa, a.Temp1Float32), matScale(sb, b.Temp1Float32))
	sigma = matAdd(sigma, matScale(outer(dax, day, daz), a.Temp1Float32))
	sigma = matAdd(sigma, matScale(outer(dbx, dby, dbz), b.Temp1Float32))
	sigma = matScale(sigma, invTotal)

	vals, vecs := eigenDecomp(sigma)
	o.ScaleX = safeLogSqrt(vals[0])
	o.ScaleY = safeLogSqrt(vals[1])
	o.ScaleZ = safeLogSqrt(vals[2])

	qw, qx, qy, qz := matToQuat(vecs)
	o.RotationW, o.RotationX, o.RotationY, o.RotationZ = cmn.NormalizeRotationsF32Uint8(qw, qx, qy, qz)

	o.Temp1Float32 = 0
	o.Temp2Float32 = 0
	o.Temp3Float32 = 0
	o.Temp4Float32 = 0
	o.SH45 = nil
	return o
}

func processChunk(chunk []*SplatData, gridStep float32) []*SplatData {
	n := len(chunk)
	if n < 2 {
		return chunk
	}

	invGrid := 1.0 / gridStep
	fh := newFlatHash(n >> 2)
	cellBuf := make([]cellInfo, n)

	for i := range chunk {
		s := chunk[i]
		gx := int32(math.Floor(float64(s.PositionX * invGrid)))
		gy := int32(math.Floor(float64(s.PositionY * invGrid)))
		gz := int32(math.Floor(float64(s.PositionZ * invGrid)))
		cellBuf[i] = cellInfo{gx, gy, gz, int32(i)}
		fh.add(gridKey{gx, gy, gz}, int32(i))
	}

	used := make([]uint8, n)
	out := make([]*SplatData, 0, n>>1)

	for i := range chunk {
		if used[i] != 0 {
			continue
		}
		s := chunk[i]
		ci := cellBuf[i]
		bestIdx := int32(-1)
		bestScore := float32(-1000000)

		for dz := int32(-1); dz <= 1; dz++ {
			for dy := int32(-1); dy <= 1; dy++ {
				for dx := int32(-1); dx <= 1; dx++ {
					vals := fh.get(gridKey{ci.gx + dx, ci.gy + dy, ci.gz + dz})
					for _, j := range vals {
						if j == int32(i) || used[j] != 0 {
							continue
						}
						ratio := s.Temp2Float32 / (chunk[j].Temp2Float32 + 1e-6)
						if ratio < scaleMinRatio || ratio > scaleMaxRatio {
							continue
						}
						score := similarity(s, chunk[j])
						if score > bestScore {
							bestScore, bestIdx = score, j
						}
					}
				}
			}
		}

		if bestIdx >= 0 && bestScore > baseMergeThreshold {
			out = append(out, mergeTwo(s, chunk[bestIdx]))
			used[i] = 1
			used[bestIdx] = 1
		} else {
			out = append(out, s)
		}
	}
	return out
}

func safeLogSqrt(v float32) float32 {
	if v < 1e-9 {
		v = 1e-9
	}
	return float32(math.Log(math.Sqrt(float64(v))))
}
