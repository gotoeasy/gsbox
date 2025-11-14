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
	var width int
	if meta.ShN != nil {
		centroids, width = webpRgbaWidth(filepath.Join(dir, meta.ShN.Files[0]))
		labels = webpRgba(filepath.Join(dir, meta.ShN.Files[1]))
	}

	count := meta.Count
	datas := make([]*SplatData, count)
	shDegree := uint8(0)
	if meta.ShN != nil {
		shDegree = 3
	}

	for i := range count {
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
			label := int(labels[i*4+0]) + (int(labels[i*4+1]) << 8)
			col := (label & 63) // 同 (n % 64)
			row := label >> 6   // 同 Math.floor(n / 64)
			offset := row*width + col*15

			sh1 := make([]float32, 9)
			sh2 := make([]float32, 15)
			sh3 := make([]float32, 21)
			for d := range 3 {
				for k := range 3 {
					sh1[k*3+d] = meta.ShN.Codebook[centroids[(offset+k)*4+d]]
				}
				for k := range 5 {
					sh2[k*3+d] = meta.ShN.Codebook[centroids[(offset+3+k)*4+d]]
				}
				for k := range 7 {
					sh3[k*3+d] = meta.ShN.Codebook[centroids[(offset+8+k)*4+d]]
				}
			}
			var shs []uint8
			for _, val := range sh1 {
				shs = append(shs, cmn.EncodeSplatSH(float64(val)))
			}
			for _, val := range sh2 {
				shs = append(shs, cmn.EncodeSplatSH(float64(val)))
			}
			for _, val := range sh3 {
				shs = append(shs, cmn.EncodeSplatSH(float64(val)))
			}

			switch shDegree {
			case 3:
				splatData.SH2 = shs[:24]
				splatData.SH3 = shs[24:]
			case 2:
				splatData.SH2 = shs[:24]
			case 1:
				splatData.SH1 = shs[:9]
			}
		}

		datas[i] = splatData
	}

	header := &SogHeader{}
	header.Version = 2
	header.Count = count
	header.ShDegree = shDegree
	inputSogHeader = header
	return datas, header
}
