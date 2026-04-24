package gsplat

import (
	"gsbox/cmn"
	"math"
	"sort"
)

// 高斯简化
func Simplify(splats []*SplatData) []*SplatData {
	n := len(splats)
	if n < 2 {
		return splats
	}

	type gridKey struct{ x, y, z int }

	var avgScale float32
	for i := range splats {
		alpha := cmn.EncodeSplatOpacityF32(splats[i].ColorA)

		sx := float32(math.Exp(float64(splats[i].ScaleX)))
		sy := float32(math.Exp(float64(splats[i].ScaleY)))
		sz := float32(math.Exp(float64(splats[i].ScaleZ)))

		volume := sx * sy * sz
		splats[i].TempSi = volume * alpha

		avgScale += (sx + sy + sz) / 3.0
	}
	avgScale /= float32(n)

	sort.Slice(splats, func(i, j int) bool {
		return splats[i].TempSi > splats[j].TempSi
	})

	gridSize := avgScale * 2.0
	if gridSize < 1e-3 {
		gridSize = 0.5
	}

	cellMap := make(map[gridKey][]int)
	for idx := range splats {
		s := splats[idx]
		gx := int(math.Floor(float64(s.PositionX / gridSize)))
		gy := int(math.Floor(float64(s.PositionY / gridSize)))
		gz := int(math.Floor(float64(s.PositionZ / gridSize)))
		key := gridKey{gx, gy, gz}
		cellMap[key] = append(cellMap[key], idx)
	}

	used := make([]bool, n)
	result := make([]*SplatData, 0, n/2)
	const simThreshold = float32(-3.0)

	for i := range splats {
		if used[i] {
			continue
		}

		s := splats[i]
		gx := int(math.Floor(float64(s.PositionX / gridSize)))
		gy := int(math.Floor(float64(s.PositionY / gridSize)))
		gz := int(math.Floor(float64(s.PositionZ / gridSize)))

		bestIdx := -1
		bestScore := float32(-1e9)

		for dz := -1; dz <= 1; dz++ {
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					key := gridKey{gx + dx, gy + dy, gz + dz}
					if list, ok := cellMap[key]; ok {
						for _, j := range list {
							if j == i || used[j] {
								continue
							}

							sa := math.Exp(float64(s.ScaleX))
							sb := math.Exp(float64(splats[j].ScaleX))
							ratio := sa / (sb + 1e-6)
							if ratio > 3.0 || ratio < 0.5 {
								continue
							}

							score := similarity(s, splats[j])
							if score > bestScore {
								bestScore = score
								bestIdx = j
							}
						}
					}
				}
			}
		}

		if bestIdx != -1 && bestScore > simThreshold {
			result = append(result, mergeSplats(s, splats[bestIdx]))
			used[i] = true
			used[bestIdx] = true
		} else {
			result = append(result, s)
		}
	}

	return result
}

// 计算匹配度
func similarity(a, b *SplatData) float32 {
	dx := a.PositionX - b.PositionX
	dy := a.PositionY - b.PositionY
	dz := a.PositionZ - b.PositionZ
	distSq := dx*dx + dy*dy + dz*dz

	dr := (float32(a.ColorR) - float32(b.ColorR)) / 255.0
	dg := (float32(a.ColorG) - float32(b.ColorG)) / 255.0
	db := (float32(a.ColorB) - float32(b.ColorB)) / 255.0
	colorSq := dr*dr + dg*dg + db*db

	sa := float32(math.Exp(float64(a.ScaleX)))
	sb := float32(math.Exp(float64(b.ScaleX)))
	scaleDiff := float32(math.Abs(float64(sa-sb))) / (sa + sb + 1e-6)

	return -distSq*1.0 - colorSq*0.5 - scaleDiff*2.0
}

// 高斯合并
func mergeSplats(a, b *SplatData) *SplatData {
	wA := a.TempSi
	wB := b.TempSi
	totalW := wA + wB
	if totalW < 1e-6 {
		return a
	}

	mx := (a.PositionX*wA + b.PositionX*wB) / totalW
	my := (a.PositionY*wA + b.PositionY*wB) / totalW
	mz := (a.PositionZ*wA + b.PositionZ*wB) / totalW

	cr := (float32(a.ColorR)*wA + float32(b.ColorR)*wB) / totalW
	cg := (float32(a.ColorG)*wA + float32(b.ColorG)*wB) / totalW
	cb := (float32(a.ColorB)*wA + float32(b.ColorB)*wB) / totalW
	ca := (float32(a.ColorA)*wA + float32(b.ColorA)*wB) / totalW

	sigmaA := buildSigma(a)
	sigmaB := buildSigma(b)

	dax := a.PositionX - mx
	day := a.PositionY - my
	daz := a.PositionZ - mz

	dbx := b.PositionX - mx
	dby := b.PositionY - my
	dbz := b.PositionZ - mz

	sigma := matAdd(matScale(sigmaA, wA), matScale(sigmaB, wB))
	sigma = matAdd(sigma, matScale(outer(dax, day, daz), wA))
	sigma = matAdd(sigma, matScale(outer(dbx, dby, dbz), wB))
	sigma = matScale(sigma, 1.0/totalW)

	vals, vecs := eigenDecomp(sigma)

	sx := float32(math.Sqrt(float64(vals[0])))
	sy := float32(math.Sqrt(float64(vals[1])))
	sz := float32(math.Sqrt(float64(vals[2])))

	logSx := float32(math.Log(float64(sx) + 1e-9))
	logSy := float32(math.Log(float64(sy) + 1e-9))
	logSz := float32(math.Log(float64(sz) + 1e-9))

	qw, qx, qy, qz := matToQuat(vecs)
	rw, rx, ry, rz := cmn.NormalizeRotationsF32Uint8(qw, qx, qy, qz)

	return &SplatData{
		PositionX: mx, PositionY: my, PositionZ: mz,
		ScaleX: logSx, ScaleY: logSy, ScaleZ: logSz,
		ColorR: uint8(cr), ColorG: uint8(cg), ColorB: uint8(cb), ColorA: uint8(ca),
		RotationW: rw, RotationX: rx, RotationY: ry, RotationZ: rz,
	}
}

func quatToMat3(w, x, y, z float32) [3][3]float32 {
	return [3][3]float32{{1 - 2*(y*y+z*z), 2 * (x*y - z*w), 2 * (x*z + y*w)}, {2 * (x*y + z*w), 1 - 2*(x*x+z*z), 2 * (y*z - x*w)}, {2 * (x*z - y*w), 2 * (y*z + x*w), 1 - 2*(x*x+y*y)}}
}

func matMul(a, b [3][3]float32) [3][3]float32 {
	var r [3][3]float32
	for i := range 3 {
		for j := range 3 {
			r[i][j] = a[i][0]*b[0][j] + a[i][1]*b[1][j] + a[i][2]*b[2][j]
		}
	}
	return r
}

func matT(m [3][3]float32) [3][3]float32 {
	return [3][3]float32{{m[0][0], m[1][0], m[2][0]}, {m[0][1], m[1][1], m[2][1]}, {m[0][2], m[1][2], m[2][2]}}
}

func matAdd(a, b [3][3]float32) [3][3]float32 {
	var r [3][3]float32
	for i := range 3 {
		for j := range 3 {
			r[i][j] = a[i][j] + b[i][j]
		}
	}
	return r
}

func matScale(a [3][3]float32, s float32) [3][3]float32 {
	var r [3][3]float32
	for i := range 3 {
		for j := range 3 {
			r[i][j] = a[i][j] * s
		}
	}
	return r
}

func outer(x, y, z float32) [3][3]float32 {
	return [3][3]float32{{x * x, x * y, x * z}, {y * x, y * y, y * z}, {z * x, z * y, z * z}}
}

func buildSigma(s *SplatData) [3][3]float32 {
	qw, qx, qy, qz := cmn.NormalizeRotationsUint8F32(s.RotationW, s.RotationX, s.RotationY, s.RotationZ)
	R := quatToMat3(qw, qx, qy, qz)

	sx := float32(math.Exp(float64(s.ScaleX)))
	sy := float32(math.Exp(float64(s.ScaleY)))
	sz := float32(math.Exp(float64(s.ScaleZ)))

	D := [3][3]float32{{sx * sx, 0, 0}, {0, sy * sy, 0}, {0, 0, sz * sz}}

	return matMul(matMul(R, D), matT(R))
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
