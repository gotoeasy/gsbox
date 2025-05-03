package gsplat

import (
	"gsbox/cmn"
	"math"
)

var kSqrt03_02 = math.Sqrt(3.0 / 2.0)
var kSqrt01_03 = math.Sqrt(1.0 / 3.0)
var kSqrt02_03 = math.Sqrt(2.0 / 3.0)
var kSqrt04_03 = math.Sqrt(4.0 / 3.0)
var kSqrt01_04 = math.Sqrt(1.0 / 4.0)
var kSqrt03_04 = math.Sqrt(3.0 / 4.0)
var kSqrt01_05 = math.Sqrt(1.0 / 5.0)
var kSqrt03_05 = math.Sqrt(3.0 / 5.0)
var kSqrt06_05 = math.Sqrt(6.0 / 5.0)
var kSqrt08_05 = math.Sqrt(8.0 / 5.0)
var kSqrt09_05 = math.Sqrt(9.0 / 5.0)
var kSqrt01_06 = math.Sqrt(1.0 / 6.0)
var kSqrt05_06 = math.Sqrt(5.0 / 6.0)
var kSqrt03_08 = math.Sqrt(3.0 / 8.0)
var kSqrt05_08 = math.Sqrt(5.0 / 8.0)
var kSqrt09_08 = math.Sqrt(9.0 / 8.0)
var kSqrt05_09 = math.Sqrt(5.0 / 9.0)
var kSqrt08_09 = math.Sqrt(8.0 / 9.0)
var kSqrt01_10 = math.Sqrt(1.0 / 10.0)
var kSqrt03_10 = math.Sqrt(3.0 / 10.0)
var kSqrt01_12 = math.Sqrt(1.0 / 12.0)
var kSqrt04_15 = math.Sqrt(4.0 / 15.0)
var kSqrt01_16 = math.Sqrt(1.0 / 16.0)
var kSqrt15_16 = math.Sqrt(15.0 / 16.0)
var kSqrt01_18 = math.Sqrt(1.0 / 18.0)
var kSqrt01_60 = math.Sqrt(1.0 / 60.0)

type SHRotation struct {
	sh1 [3][]float64
	sh2 [5][]float64
	sh3 [7][]float64
}

// from https://github.com/playcanvas/supersplat/blob/main/src/sh-utils.ts
func NewSHRotation(q *Quaternion) *SHRotation {
	rot := quaternionToMat3Array(q)

	// band 1
	sh1 := [3][]float64{
		{rot[4], -rot[7], rot[1]},
		{-rot[5], rot[8], -rot[2]},
		{rot[3], -rot[6], rot[0]},
	}

	// band 2
	sh2 := [5][]float64{{
		kSqrt01_04 * ((sh1[2][2]*sh1[0][0] + sh1[2][0]*sh1[0][2]) + (sh1[0][2]*sh1[2][0] + sh1[0][0]*sh1[2][2])),
		(sh1[2][1]*sh1[0][0] + sh1[0][1]*sh1[2][0]),
		kSqrt03_04 * (sh1[2][1]*sh1[0][1] + sh1[0][1]*sh1[2][1]),
		(sh1[2][1]*sh1[0][2] + sh1[0][1]*sh1[2][2]),
		kSqrt01_04 * ((sh1[2][2]*sh1[0][2] - sh1[2][0]*sh1[0][0]) + (sh1[0][2]*sh1[2][2] - sh1[0][0]*sh1[2][0])),
	}, {
		kSqrt01_04 * ((sh1[1][2]*sh1[0][0] + sh1[1][0]*sh1[0][2]) + (sh1[0][2]*sh1[1][0] + sh1[0][0]*sh1[1][2])),
		sh1[1][1]*sh1[0][0] + sh1[0][1]*sh1[1][0],
		kSqrt03_04 * (sh1[1][1]*sh1[0][1] + sh1[0][1]*sh1[1][1]),
		sh1[1][1]*sh1[0][2] + sh1[0][1]*sh1[1][2],
		kSqrt01_04 * ((sh1[1][2]*sh1[0][2] - sh1[1][0]*sh1[0][0]) + (sh1[0][2]*sh1[1][2] - sh1[0][0]*sh1[1][0])),
	}, {
		kSqrt01_03*(sh1[1][2]*sh1[1][0]+sh1[1][0]*sh1[1][2]) - kSqrt01_12*((sh1[2][2]*sh1[2][0]+sh1[2][0]*sh1[2][2])+(sh1[0][2]*sh1[0][0]+sh1[0][0]*sh1[0][2])),
		kSqrt04_03*sh1[1][1]*sh1[1][0] - kSqrt01_03*(sh1[2][1]*sh1[2][0]+sh1[0][1]*sh1[0][0]),
		sh1[1][1]*sh1[1][1] - kSqrt01_04*(sh1[2][1]*sh1[2][1]+sh1[0][1]*sh1[0][1]),
		kSqrt04_03*sh1[1][1]*sh1[1][2] - kSqrt01_03*(sh1[2][1]*sh1[2][2]+sh1[0][1]*sh1[0][2]),
		kSqrt01_03*(sh1[1][2]*sh1[1][2]-sh1[1][0]*sh1[1][0]) - kSqrt01_12*((sh1[2][2]*sh1[2][2]-sh1[2][0]*sh1[2][0])+(sh1[0][2]*sh1[0][2]-sh1[0][0]*sh1[0][0])),
	}, {
		kSqrt01_04 * ((sh1[1][2]*sh1[2][0] + sh1[1][0]*sh1[2][2]) + (sh1[2][2]*sh1[1][0] + sh1[2][0]*sh1[1][2])),
		sh1[1][1]*sh1[2][0] + sh1[2][1]*sh1[1][0],
		kSqrt03_04 * (sh1[1][1]*sh1[2][1] + sh1[2][1]*sh1[1][1]),
		sh1[1][1]*sh1[2][2] + sh1[2][1]*sh1[1][2],
		kSqrt01_04 * ((sh1[1][2]*sh1[2][2] - sh1[1][0]*sh1[2][0]) + (sh1[2][2]*sh1[1][2] - sh1[2][0]*sh1[1][0])),
	}, {
		kSqrt01_04 * ((sh1[2][2]*sh1[2][0] + sh1[2][0]*sh1[2][2]) - (sh1[0][2]*sh1[0][0] + sh1[0][0]*sh1[0][2])),
		(sh1[2][1]*sh1[2][0] - sh1[0][1]*sh1[0][0]),
		kSqrt03_04 * (sh1[2][1]*sh1[2][1] - sh1[0][1]*sh1[0][1]),
		(sh1[2][1]*sh1[2][2] - sh1[0][1]*sh1[0][2]),
		kSqrt01_04 * ((sh1[2][2]*sh1[2][2] - sh1[2][0]*sh1[2][0]) - (sh1[0][2]*sh1[0][2] - sh1[0][0]*sh1[0][0])),
	}}

	// band 3
	sh3 := [7][]float64{{
		kSqrt01_04 * ((sh1[2][2]*sh2[0][0] + sh1[2][0]*sh2[0][4]) + (sh1[0][2]*sh2[4][0] + sh1[0][0]*sh2[4][4])),
		kSqrt03_02 * (sh1[2][1]*sh2[0][0] + sh1[0][1]*sh2[4][0]),
		kSqrt15_16 * (sh1[2][1]*sh2[0][1] + sh1[0][1]*sh2[4][1]),
		kSqrt05_06 * (sh1[2][1]*sh2[0][2] + sh1[0][1]*sh2[4][2]),
		kSqrt15_16 * (sh1[2][1]*sh2[0][3] + sh1[0][1]*sh2[4][3]),
		kSqrt03_02 * (sh1[2][1]*sh2[0][4] + sh1[0][1]*sh2[4][4]),
		kSqrt01_04 * ((sh1[2][2]*sh2[0][4] - sh1[2][0]*sh2[0][0]) + (sh1[0][2]*sh2[4][4] - sh1[0][0]*sh2[4][0])),
	}, {
		kSqrt01_06*(sh1[1][2]*sh2[0][0]+sh1[1][0]*sh2[0][4]) + kSqrt01_06*((sh1[2][2]*sh2[1][0]+sh1[2][0]*sh2[1][4])+(sh1[0][2]*sh2[3][0]+sh1[0][0]*sh2[3][4])),
		sh1[1][1]*sh2[0][0] + (sh1[2][1]*sh2[1][0] + sh1[0][1]*sh2[3][0]),
		kSqrt05_08*sh1[1][1]*sh2[0][1] + kSqrt05_08*(sh1[2][1]*sh2[1][1]+sh1[0][1]*sh2[3][1]),
		kSqrt05_09*sh1[1][1]*sh2[0][2] + kSqrt05_09*(sh1[2][1]*sh2[1][2]+sh1[0][1]*sh2[3][2]),
		kSqrt05_08*sh1[1][1]*sh2[0][3] + kSqrt05_08*(sh1[2][1]*sh2[1][3]+sh1[0][1]*sh2[3][3]),
		sh1[1][1]*sh2[0][4] + (sh1[2][1]*sh2[1][4] + sh1[0][1]*sh2[3][4]),
		kSqrt01_06*(sh1[1][2]*sh2[0][4]-sh1[1][0]*sh2[0][0]) + kSqrt01_06*((sh1[2][2]*sh2[1][4]-sh1[2][0]*sh2[1][0])+(sh1[0][2]*sh2[3][4]-sh1[0][0]*sh2[3][0])),
	}, {
		kSqrt04_15*(sh1[1][2]*sh2[1][0]+sh1[1][0]*sh2[1][4]) + kSqrt01_05*(sh1[0][2]*sh2[2][0]+sh1[0][0]*sh2[2][4]) - kSqrt01_60*((sh1[2][2]*sh2[0][0]+sh1[2][0]*sh2[0][4])-(sh1[0][2]*sh2[4][0]+sh1[0][0]*sh2[4][4])),
		kSqrt08_05*sh1[1][1]*sh2[1][0] + kSqrt06_05*sh1[0][1]*sh2[2][0] - kSqrt01_10*(sh1[2][1]*sh2[0][0]-sh1[0][1]*sh2[4][0]),
		sh1[1][1]*sh2[1][1] + kSqrt03_04*sh1[0][1]*sh2[2][1] - kSqrt01_16*(sh1[2][1]*sh2[0][1]-sh1[0][1]*sh2[4][1]),
		kSqrt08_09*sh1[1][1]*sh2[1][2] + kSqrt02_03*sh1[0][1]*sh2[2][2] - kSqrt01_18*(sh1[2][1]*sh2[0][2]-sh1[0][1]*sh2[4][2]),
		sh1[1][1]*sh2[1][3] + kSqrt03_04*sh1[0][1]*sh2[2][3] - kSqrt01_16*(sh1[2][1]*sh2[0][3]-sh1[0][1]*sh2[4][3]),
		kSqrt08_05*sh1[1][1]*sh2[1][4] + kSqrt06_05*sh1[0][1]*sh2[2][4] - kSqrt01_10*(sh1[2][1]*sh2[0][4]-sh1[0][1]*sh2[4][4]),
		kSqrt04_15*(sh1[1][2]*sh2[1][4]-sh1[1][0]*sh2[1][0]) + kSqrt01_05*(sh1[0][2]*sh2[2][4]-sh1[0][0]*sh2[2][0]) - kSqrt01_60*((sh1[2][2]*sh2[0][4]-sh1[2][0]*sh2[0][0])-(sh1[0][2]*sh2[4][4]-sh1[0][0]*sh2[4][0])),
	}, {
		kSqrt03_10*(sh1[1][2]*sh2[2][0]+sh1[1][0]*sh2[2][4]) - kSqrt01_10*((sh1[2][2]*sh2[3][0]+sh1[2][0]*sh2[3][4])+(sh1[0][2]*sh2[1][0]+sh1[0][0]*sh2[1][4])),
		kSqrt09_05*sh1[1][1]*sh2[2][0] - kSqrt03_05*(sh1[2][1]*sh2[3][0]+sh1[0][1]*sh2[1][0]),
		kSqrt09_08*sh1[1][1]*sh2[2][1] - kSqrt03_08*(sh1[2][1]*sh2[3][1]+sh1[0][1]*sh2[1][1]),
		sh1[1][1]*sh2[2][2] - kSqrt01_03*(sh1[2][1]*sh2[3][2]+sh1[0][1]*sh2[1][2]),
		kSqrt09_08*sh1[1][1]*sh2[2][3] - kSqrt03_08*(sh1[2][1]*sh2[3][3]+sh1[0][1]*sh2[1][3]),
		kSqrt09_05*sh1[1][1]*sh2[2][4] - kSqrt03_05*(sh1[2][1]*sh2[3][4]+sh1[0][1]*sh2[1][4]),
		kSqrt03_10*(sh1[1][2]*sh2[2][4]-sh1[1][0]*sh2[2][0]) - kSqrt01_10*((sh1[2][2]*sh2[3][4]-sh1[2][0]*sh2[3][0])+(sh1[0][2]*sh2[1][4]-sh1[0][0]*sh2[1][0])),
	}, {
		kSqrt04_15*(sh1[1][2]*sh2[3][0]+sh1[1][0]*sh2[3][4]) + kSqrt01_05*(sh1[2][2]*sh2[2][0]+sh1[2][0]*sh2[2][4]) - kSqrt01_60*((sh1[2][2]*sh2[4][0]+sh1[2][0]*sh2[4][4])+(sh1[0][2]*sh2[0][0]+sh1[0][0]*sh2[0][4])),
		kSqrt08_05*sh1[1][1]*sh2[3][0] + kSqrt06_05*sh1[2][1]*sh2[2][0] - kSqrt01_10*(sh1[2][1]*sh2[4][0]+sh1[0][1]*sh2[0][0]),
		sh1[1][1]*sh2[3][1] + kSqrt03_04*sh1[2][1]*sh2[2][1] - kSqrt01_16*(sh1[2][1]*sh2[4][1]+sh1[0][1]*sh2[0][1]),
		kSqrt08_09*sh1[1][1]*sh2[3][2] + kSqrt02_03*sh1[2][1]*sh2[2][2] - kSqrt01_18*(sh1[2][1]*sh2[4][2]+sh1[0][1]*sh2[0][2]),
		sh1[1][1]*sh2[3][3] + kSqrt03_04*sh1[2][1]*sh2[2][3] - kSqrt01_16*(sh1[2][1]*sh2[4][3]+sh1[0][1]*sh2[0][3]),
		kSqrt08_05*sh1[1][1]*sh2[3][4] + kSqrt06_05*sh1[2][1]*sh2[2][4] - kSqrt01_10*(sh1[2][1]*sh2[4][4]+sh1[0][1]*sh2[0][4]),
		kSqrt04_15*(sh1[1][2]*sh2[3][4]-sh1[1][0]*sh2[3][0]) + kSqrt01_05*(sh1[2][2]*sh2[2][4]-sh1[2][0]*sh2[2][0]) - kSqrt01_60*((sh1[2][2]*sh2[4][4]-sh1[2][0]*sh2[4][0])+(sh1[0][2]*sh2[0][4]-sh1[0][0]*sh2[0][0])),
	}, {
		kSqrt01_06*(sh1[1][2]*sh2[4][0]+sh1[1][0]*sh2[4][4]) + kSqrt01_06*((sh1[2][2]*sh2[3][0]+sh1[2][0]*sh2[3][4])-(sh1[0][2]*sh2[1][0]+sh1[0][0]*sh2[1][4])),
		sh1[1][1]*sh2[4][0] + (sh1[2][1]*sh2[3][0] - sh1[0][1]*sh2[1][0]),
		kSqrt05_08*sh1[1][1]*sh2[4][1] + kSqrt05_08*(sh1[2][1]*sh2[3][1]-sh1[0][1]*sh2[1][1]),
		kSqrt05_09*sh1[1][1]*sh2[4][2] + kSqrt05_09*(sh1[2][1]*sh2[3][2]-sh1[0][1]*sh2[1][2]),
		kSqrt05_08*sh1[1][1]*sh2[4][3] + kSqrt05_08*(sh1[2][1]*sh2[3][3]-sh1[0][1]*sh2[1][3]),
		sh1[1][1]*sh2[4][4] + (sh1[2][1]*sh2[3][4] - sh1[0][1]*sh2[1][4]),
		kSqrt01_06*(sh1[1][2]*sh2[4][4]-sh1[1][0]*sh2[4][0]) + kSqrt01_06*((sh1[2][2]*sh2[3][4]-sh1[2][0]*sh2[3][0])-(sh1[0][2]*sh2[1][4]-sh1[0][0]*sh2[1][0])),
	}, {
		kSqrt01_04 * ((sh1[2][2]*sh2[4][0] + sh1[2][0]*sh2[4][4]) - (sh1[0][2]*sh2[0][0] + sh1[0][0]*sh2[0][4])),
		kSqrt03_02 * (sh1[2][1]*sh2[4][0] - sh1[0][1]*sh2[0][0]),
		kSqrt15_16 * (sh1[2][1]*sh2[4][1] - sh1[0][1]*sh2[0][1]),
		kSqrt05_06 * (sh1[2][1]*sh2[4][2] - sh1[0][1]*sh2[0][2]),
		kSqrt15_16 * (sh1[2][1]*sh2[4][3] - sh1[0][1]*sh2[0][3]),
		kSqrt03_02 * (sh1[2][1]*sh2[4][4] - sh1[0][1]*sh2[0][4]),
		kSqrt01_04 * ((sh1[2][2]*sh2[4][4] - sh1[2][0]*sh2[4][0]) - (sh1[0][2]*sh2[0][4] - sh1[0][0]*sh2[0][0])),
	}}

	return &SHRotation{
		sh1: sh1,
		sh2: sh2,
		sh3: sh3,
	}
}

// rotate spherical harmonic coefficients, up to band 3
func (s *SHRotation) Apply(result []float32) {
	src := make([]float32, len(result))
	copy(src, result)

	// band 1
	if len(result) < 3 {
		return
	}
	result[0] = dp(3, 0, src, s.sh1[0])
	result[1] = dp(3, 0, src, s.sh1[1])
	result[2] = dp(3, 0, src, s.sh1[2])

	// band 2
	if len(result) < 8 {
		return
	}
	result[3] = dp(5, 3, src, s.sh2[0])
	result[4] = dp(5, 3, src, s.sh2[1])
	result[5] = dp(5, 3, src, s.sh2[2])
	result[6] = dp(5, 3, src, s.sh2[3])
	result[7] = dp(5, 3, src, s.sh2[4])

	// band 3
	if len(result) < 15 {
		return
	}
	result[8] = dp(7, 8, src, s.sh3[0])
	result[9] = dp(7, 8, src, s.sh3[1])
	result[10] = dp(7, 8, src, s.sh3[2])
	result[11] = dp(7, 8, src, s.sh3[3])
	result[12] = dp(7, 8, src, s.sh3[4])
	result[13] = dp(7, 8, src, s.sh3[5])
	result[14] = dp(7, 8, src, s.sh3[6])
}

func dp(n int, start int, a []float32, b []float64) float32 {
	sum := 0.0
	for i := range n {
		sum += float64(a[start+i]) * b[i]
	}
	return cmn.ClipFloat32(sum)
}

func quaternionToMat3Array(q *Quaternion) [9]float64 {
	qx, qy, qz, qw := q.X, q.Y, q.Z, q.W

	x2 := qx + qx
	y2 := qy + qy
	z2 := qz + qz
	xx := qx * x2
	xy := qx * y2
	xz := qx * z2
	yy := qy * y2
	yz := qy * z2
	zz := qz * z2
	wx := qw * x2
	wy := qw * y2
	wz := qw * z2

	m := [9]float64{}

	m[0] = (1 - (yy + zz))
	m[1] = (xy + wz)
	m[2] = (xz - wy)

	m[3] = (xy - wz)
	m[4] = (1 - (xx + zz))
	m[5] = (yz + wx)

	m[6] = (xz + wy)
	m[7] = (yz - wx)
	m[8] = (1 - (xx + yy))

	return m
}
