package gsplat

import (
	"gsbox/cmn"
	"math"
	"path/filepath"
)

func ReadSogV2(meta *SogMeta, dir string) ([]*SplatData, *SogHeader) {
	meansl := webpRgba(filepath.Join(dir, meta.Means.Files[0]))
	meansu := webpRgba(filepath.Join(dir, meta.Means.Files[1]))
	scales := webpRgba(filepath.Join(dir, meta.Scales.Files[0]))
	quats := webpRgba(filepath.Join(dir, meta.Quats.Files[0]))
	sh0 := webpRgba(filepath.Join(dir, meta.Sh0.Files[0]))
	var centroids []byte
	var labels []byte
	var width int // 96|512|960，不同宽度对应不同级别
	var paletteSize int
	if meta.ShN != nil {
		centroids, width = webpRgbaWidth(filepath.Join(dir, meta.ShN.Files[0]))
		labels = webpRgba(filepath.Join(dir, meta.ShN.Files[1]))
		if meta.ShN.Count >= 0 {
			paletteSize = meta.ShN.Count // v2 新设计新添字段
		} else {
			paletteSize = 65536 // v2 旧设计无调色板数量字段，按最大值记
		}
	}

	count := meta.Count
	datas := make([]*SplatData, count)
	shDegree := uint8(0)
	if meta.ShN != nil {
		if meta.ShN.Count > 0 && meta.ShN.Bands > 0 {
			shDegree = uint8(meta.ShN.Bands) // 补丁，版本2的早期设计无Count、Bands字段
		} else {
			shDegree = 3 // 早期默认都是级别3
		}
	}

	for i := range count {
		OnProgress(PhaseRead, i, count)
		splatData := &SplatData{}

		fx := float32((uint16(meansu[i*4+0])<<8)|uint16(meansl[i*4+0])) / 65535.0
		fy := float32((uint16(meansu[i*4+1])<<8)|uint16(meansl[i*4+1])) / 65535.0
		fz := float32((uint16(meansu[i*4+2])<<8)|uint16(meansl[i*4+2])) / 65535.0
		x := meta.Means.Mins[0] + (meta.Means.Maxs[0]-meta.Means.Mins[0])*fx
		y := meta.Means.Mins[1] + (meta.Means.Maxs[1]-meta.Means.Mins[1])*fy
		z := meta.Means.Mins[2] + (meta.Means.Maxs[2]-meta.Means.Mins[2])*fz
		if x < 0 {
			x = -cmn.ClipFloat32(math.Exp(math.Abs(float64(x))) - 1.0)
		} else {
			x = cmn.ClipFloat32(math.Exp(math.Abs(float64(x))) - 1.0)
		}
		if y < 0 {
			y = -cmn.ClipFloat32(math.Exp(math.Abs(float64(y))) - 1.0)
		} else {
			y = cmn.ClipFloat32(math.Exp(math.Abs(float64(y))) - 1.0)
		}
		if z < 0 {
			z = -cmn.ClipFloat32(math.Exp(math.Abs(float64(z))) - 1.0)
		} else {
			z = cmn.ClipFloat32(math.Exp(math.Abs(float64(z))) - 1.0)
		}
		splatData.PositionX, splatData.PositionY, splatData.PositionZ = x, y, z

		sx := meta.Scales.Codebook[scales[i*4+0]]
		sy := meta.Scales.Codebook[scales[i*4+1]]
		sz := meta.Scales.Codebook[scales[i*4+2]]
		splatData.ScaleX, splatData.ScaleY, splatData.ScaleZ = sx, sy, sz

		splatData.RotationW, splatData.RotationX, splatData.RotationY, splatData.RotationZ = cmn.SogDecodeRotations(quats[i*4+0], quats[i*4+1], quats[i*4+2], uint8(quats[i*4+3]))

		r := float64(meta.Sh0.Codebook[sh0[i*4+0]])
		g := float64(meta.Sh0.Codebook[sh0[i*4+1]])
		b := float64(meta.Sh0.Codebook[sh0[i*4+2]])
		a := sh0[i*4+3]
		splatData.ColorR, splatData.ColorG, splatData.ColorB, splatData.ColorA = cmn.EncodeSplatColor(r), cmn.EncodeSplatColor(g), cmn.EncodeSplatColor(b), a

		if shDegree > 0 {
			shDims := []int{0, 3, 8, 15}
			shDim := shDims[shDegree]
			label := int(labels[i*4+0]) | (int(labels[i*4+1]) << 8)
			col := (label & 63)             // 同 (n % 64)
			row := label >> 6               // 同 Math.floor(n / 64)
			offset := row*width + col*shDim // 像素单位的偏移量
			f32shs := make([]float32, 45)   // rgb0,rgb1 ... rgb14

			// 避免数据问题越界
			if label < paletteSize {
				for k := range shDim {
					f32shs[k*3+0] = meta.ShN.Codebook[centroids[(offset+k)*4+0]]
					f32shs[k*3+1] = meta.ShN.Codebook[centroids[(offset+k)*4+1]]
					f32shs[k*3+2] = meta.ShN.Codebook[centroids[(offset+k)*4+2]]
				}
			}

			shs := make([]uint8, 45)
			for n, v := range f32shs {
				shs[n] = cmn.EncodeSplatSH(float64(v))
			}

			splatData.SH45 = shs
			splatData.PaletteIdx = uint16(label)
		}

		datas[i] = splatData
	}

	header := &SogHeader{}
	header.Version = 2
	header.Count = count
	header.ShDegree = shDegree
	header.PaletteSize = paletteSize

	OnProgress(PhaseRead, 100, 100)
	return datas, header
}
