package gsplat

import (
	"gsbox/cmn"
	"math"
	"path/filepath"
)

const SQRT2 = float32(1.4142135623730951) // math.Sqrt(2.0)

func ReadSogV1(meta *SogMeta, dir string) ([]*SplatData, int) {
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

	count := meta.Means.Shape[0]
	datas := make([]*SplatData, count)
	shDegree := 0
	if meta.ShN != nil {
		switch meta.ShN.Shape[1] {
		case 45, 15:
			shDegree = 3
		case 24, 8:
			shDegree = 2
		case 9, 3:
			shDegree = 1
		}
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

		sx := float32(scales[i*4+0]) / 255.0
		sy := float32(scales[i*4+1]) / 255.0
		sz := float32(scales[i*4+2]) / 255.0
		sx = meta.Scales.Mins[0] + (meta.Scales.Maxs[0]-meta.Scales.Mins[0])*sx
		sy = meta.Scales.Mins[1] + (meta.Scales.Maxs[1]-meta.Scales.Mins[1])*sy
		sz = meta.Scales.Mins[2] + (meta.Scales.Maxs[2]-meta.Scales.Mins[2])*sz
		splatData.ScaleX, splatData.ScaleY, splatData.ScaleZ = sx, sy, sz

		r0 := (float32(quats[i*4+0])/255.0 - 0.5) * SQRT2
		r1 := (float32(quats[i*4+1])/255.0 - 0.5) * SQRT2
		r2 := (float32(quats[i*4+2])/255.0 - 0.5) * SQRT2
		ri := cmn.ClipFloat32(math.Sqrt(float64(max(0, 1.0-r0*r0-r1*r1-r2*r2))))
		idx := uint8(quats[i*4+3]) - 252
		switch idx {
		case 0:
			splatData.RotationW, splatData.RotationX, splatData.RotationY, splatData.RotationZ = cmn.NormalizeRotations4(ri, r0, r1, r2)
		case 1:
			splatData.RotationW, splatData.RotationX, splatData.RotationY, splatData.RotationZ = cmn.NormalizeRotations4(r0, ri, r1, r2)
		case 2:
			splatData.RotationW, splatData.RotationX, splatData.RotationY, splatData.RotationZ = cmn.NormalizeRotations4(r0, r1, ri, r2)
		case 3:
			splatData.RotationW, splatData.RotationX, splatData.RotationY, splatData.RotationZ = cmn.NormalizeRotations4(r0, r1, r2, ri)
		}

		r := float64(meta.Sh0.Mins[0] + (meta.Sh0.Maxs[0]-meta.Sh0.Mins[0])*(float32(sh0[i*4+0])/255.0))
		g := float64(meta.Sh0.Mins[1] + (meta.Sh0.Maxs[1]-meta.Sh0.Mins[1])*(float32(sh0[i*4+1])/255.0))
		b := float64(meta.Sh0.Mins[2] + (meta.Sh0.Maxs[2]-meta.Sh0.Mins[2])*(float32(sh0[i*4+2])/255.0))
		a := float64(meta.Sh0.Mins[3] + (meta.Sh0.Maxs[3]-meta.Sh0.Mins[3])*(float32(sh0[i*4+3])/255.0))
		splatData.ColorR, splatData.ColorG, splatData.ColorB, splatData.ColorA = cmn.EncodeSplatColor(r), cmn.EncodeSplatColor(g), cmn.EncodeSplatColor(b), cmn.EncodeSplatOpacity(a)

		if shDegree > 0 {
			label := int(labels[i*4+0]) + (int(labels[i*4+1]) << 8)
			col := (label & 63) * 15 // 同 (n % 64) * 15
			row := label >> 6        // 同 Math.floor(n / 64)
			offset := row*width + col

			sh1 := make([]float32, 9)
			sh2 := make([]float32, 15)
			sh3 := make([]float32, 21)
			for d := range 3 {
				if shDegree >= 1 {
					for k := range 3 {
						sh1[k*3+d] = ((meta.ShN.Maxs - meta.ShN.Mins) * float32(centroids[(offset+k)*4+d]) / 255.0) + meta.ShN.Mins
					}
				}
				if shDegree >= 2 {
					for k := range 5 {
						sh2[k*3+d] = ((meta.ShN.Maxs - meta.ShN.Mins) * float32(centroids[(offset+3+k)*4+d]) / 255.0) + meta.ShN.Mins
					}
				}
				if shDegree == 3 {
					for k := range 7 {
						sh3[k*3+d] = ((meta.ShN.Maxs - meta.ShN.Mins) * float32(centroids[(offset+8+k)*4+d]) / 255.0) + meta.ShN.Mins
					}
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

	return datas, shDegree
}

func webpRgba(fileWebp string) []byte {
	rgba, _ := webpRgbaWidth(fileWebp)
	return rgba
}

func webpRgbaWidth(fileWebp string) ([]byte, int) {
	webpMeansl, err := cmn.ReadFileBytes(fileWebp)
	cmn.ExitOnError(err)
	rgba, width, _, err := cmn.DecompressWebp(webpMeansl)
	cmn.ExitOnError(err)
	return rgba, width
}
