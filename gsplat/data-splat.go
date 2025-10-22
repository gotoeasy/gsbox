package gsplat

import (
	"fmt"
	"gsbox/cmn"
	"log"
	"math"
	"sort"
)

const SPLAT_DATA_SIZE = 3*4 + 3*4 + 4 + 4

type SplatData struct {
	PositionX   float32
	PositionY   float32
	PositionZ   float32
	ScaleX      float32
	ScaleY      float32
	ScaleZ      float32
	ColorR      uint8
	ColorG      uint8
	ColorB      uint8
	ColorA      uint8
	RotationW   uint8
	RotationX   uint8
	RotationY   uint8
	RotationZ   uint8
	SH1         []uint8 // sh1 only
	SH2         []uint8 // sh1 + sh2
	SH3         []uint8 // sh3 only
	IsWaterMark bool
	FlagValue   uint16
}

func TransformDatas(datas []*SplatData) []*SplatData {
	order := cmn.ToLower(Args.GetArgIgnorecase("-to", "--transform-order"))
	switch order {
	case "rts":
		transformRotateDatas(datas)
		transformTranslateDatas(datas)
		transformScaleDatas(datas)
	case "srt":
		transformScaleDatas(datas)
		transformRotateDatas(datas)
		transformTranslateDatas(datas)
	case "str":
		transformScaleDatas(datas)
		transformTranslateDatas(datas)
		transformRotateDatas(datas)
	case "trs":
		transformTranslateDatas(datas)
		transformRotateDatas(datas)
		transformScaleDatas(datas)
	case "tsr":
		transformTranslateDatas(datas)
		transformScaleDatas(datas)
		transformRotateDatas(datas)
	default:
		transformRotateDatas(datas)
		transformScaleDatas(datas)
		transformTranslateDatas(datas)
	}

	return datas
}

func transformRotateDatas(datas []*SplatData) {
	// 1, 旋转
	hasRotate, degreeX, degreeY, degreeZ := getRotateArgs()
	if hasRotate {
		qx := NewQuaternion(0, 0, 0, 1).SetFromAxisAngle(NewVector3(1, 0, 0), cmn.DegToRad(float64(degreeX)))
		qy := NewQuaternion(0, 0, 0, 1).SetFromAxisAngle(NewVector3(0, 1, 0), cmn.DegToRad(float64(degreeY)))
		qz := NewQuaternion(0, 0, 0, 1).SetFromAxisAngle(NewVector3(0, 0, 1), cmn.DegToRad(float64(degreeZ)))

		q := NewQuaternion(0, 0, 0, 1)
		if degreeX != 0 {
			q.Premultiply(qx)
		}
		if degreeY != 0 {
			q.Premultiply(qy)
		}
		if degreeZ != 0 {
			q.Premultiply(qz)
		}
		q.Normalize()
		shr := NewSHRotation(q)

		for _, data := range datas {
			data.Rotate(degreeX, degreeY, degreeZ, shr)
		}

		log.Println("[Info] (Transform) rotate in XYZ order.", "degreeX:", degreeX, ", degreeY:", degreeY, ", degreeZ:", degreeZ)
	}
}

func transformScaleDatas(datas []*SplatData) {
	// 2, 缩放
	hasScale, scale := getScaleArgs()
	if hasScale {
		for _, data := range datas {
			data.Scale(scale)
		}
		log.Println("[Info] (Transform) scaling factor:", scale)
		if scale < 0.05 {
			log.Println("[Warn] ATTENTION: VERY SMALL SCALING FACTOR MAY CAUSE PRECISION LOSS! PROCEED WITH CAUTION!")
		} else if scale > 20 {
			log.Println("[Warn] ATTENTION: VERY BIG SCALING FACTOR MAY CAUSE PRECISION LOSS! PROCEED WITH CAUTION!")
		}
	}
}

func transformTranslateDatas(datas []*SplatData) {
	// 3, 平移
	hasTranslate, tx, ty, tz := getTranslateArgs()
	if hasTranslate {
		for _, data := range datas {
			data.Translate(tx, ty, tz)
		}
		log.Println("[Info] (Transform) make translate.", "translateX:", tx, ", translateY:", ty, ", translateZ:", tz)
	}
}

func (s *SplatData) Translate(tx, ty, tz float32) {
	s.PositionX += tx
	s.PositionY += ty
	s.PositionZ += tz
}

func (s *SplatData) Scale(scale float32) {
	s.PositionX *= scale
	s.PositionY *= scale
	s.PositionZ *= scale
	s.ScaleX = cmn.DecodeSplatScale(cmn.EncodeSplatScale(s.ScaleX) * scale)
	s.ScaleY = cmn.DecodeSplatScale(cmn.EncodeSplatScale(s.ScaleY) * scale)
	s.ScaleZ = cmn.DecodeSplatScale(cmn.EncodeSplatScale(s.ScaleZ) * scale)
}

func (s *SplatData) Rotate(degreeX, degreeY, degreeZ float32, SHR *SHRotation) {

	qx := NewQuaternion(0, 0, 0, 1).SetFromAxisAngle(NewVector3(1, 0, 0), cmn.DegToRad(float64(degreeX)))
	qy := NewQuaternion(0, 0, 0, 1).SetFromAxisAngle(NewVector3(0, 1, 0), cmn.DegToRad(float64(degreeY)))
	qz := NewQuaternion(0, 0, 0, 1).SetFromAxisAngle(NewVector3(0, 0, 1), cmn.DegToRad(float64(degreeZ)))

	// rotation
	q := NewQuaternion(float64(cmn.DecodeSplatRotation(s.RotationX)), float64(cmn.DecodeSplatRotation(s.RotationY)), float64(cmn.DecodeSplatRotation(s.RotationZ)), float64(cmn.DecodeSplatRotation(s.RotationW)))
	if degreeX != 0 {
		q.Premultiply(qx)
	}
	if degreeY != 0 {
		q.Premultiply(qy)
	}
	if degreeZ != 0 {
		q.Premultiply(qz)
	}
	s.RotationW, s.RotationX, s.RotationY, s.RotationZ = cmn.NormalizeRotations(cmn.EncodeSplatRotation(q.W), cmn.EncodeSplatRotation(q.X), cmn.EncodeSplatRotation(q.Y), cmn.EncodeSplatRotation(q.Z))

	// position
	q = NewQuaternion(0, 0, 0, 1)
	if degreeX != 0 {
		q.Premultiply(qx)
	}
	if degreeY != 0 {
		q.Premultiply(qy)
	}
	if degreeZ != 0 {
		q.Premultiply(qz)
	}
	q.Normalize()
	point := NewVector3(float64(s.PositionX), float64(s.PositionY), float64(s.PositionZ))
	point.ApplyQuaternion(q)
	s.PositionX, s.PositionY, s.PositionZ = cmn.ClipFloat32(point.X), cmn.ClipFloat32(point.Y), cmn.ClipFloat32(point.Z)

	// SH
	if len(s.SH3) > 0 {
		var sh1r, sh1g, sh1b []float32
		for i := range 8 {
			sh1r = append(sh1r, cmn.DecodeSplatSH(s.SH2[i*3]))
			sh1g = append(sh1g, cmn.DecodeSplatSH(s.SH2[i*3+1]))
			sh1b = append(sh1b, cmn.DecodeSplatSH(s.SH2[i*3+2]))
		}
		for i := range 7 {
			sh1r = append(sh1r, cmn.DecodeSplatSH(s.SH3[i*3]))
			sh1g = append(sh1g, cmn.DecodeSplatSH(s.SH3[i*3+1]))
			sh1b = append(sh1b, cmn.DecodeSplatSH(s.SH3[i*3+2]))
		}
		SHR.Apply(sh1r)
		SHR.Apply(sh1g)
		SHR.Apply(sh1b)
		for i := range 8 {
			s.SH2[i*3] = cmn.EncodeSplatSH(float64(sh1r[i]))
			s.SH2[i*3+1] = cmn.EncodeSplatSH(float64(sh1g[i]))
			s.SH2[i*3+2] = cmn.EncodeSplatSH(float64(sh1b[i]))
		}
		for i := range 7 {
			s.SH3[i*3] = cmn.EncodeSplatSH(float64(sh1r[8+i]))
			s.SH3[i*3+1] = cmn.EncodeSplatSH(float64(sh1g[8+i]))
			s.SH3[i*3+2] = cmn.EncodeSplatSH(float64(sh1b[8+i]))
		}
	} else if len(s.SH2) > 0 {
		var sh1r, sh1g, sh1b []float32
		for i := range 8 {
			sh1r = append(sh1r, cmn.DecodeSplatSH(s.SH2[i*3]))
			sh1g = append(sh1g, cmn.DecodeSplatSH(s.SH2[i*3+1]))
			sh1b = append(sh1b, cmn.DecodeSplatSH(s.SH2[i*3+2]))
		}
		SHR.Apply(sh1r)
		SHR.Apply(sh1g)
		SHR.Apply(sh1b)
		for i := range 8 {
			s.SH2[i*3] = cmn.EncodeSplatSH(float64(sh1r[i]))
			s.SH2[i*3+1] = cmn.EncodeSplatSH(float64(sh1g[i]))
			s.SH2[i*3+2] = cmn.EncodeSplatSH(float64(sh1b[i]))
		}
	} else if len(s.SH1) > 0 {
		var sh1r, sh1g, sh1b []float32
		for i := range 3 {
			sh1r = append(sh1r, cmn.DecodeSplatSH(s.SH1[i*3]))
			sh1g = append(sh1g, cmn.DecodeSplatSH(s.SH1[i*3+1]))
			sh1b = append(sh1b, cmn.DecodeSplatSH(s.SH1[i*3+2]))
		}
		SHR.Apply(sh1r)
		SHR.Apply(sh1g)
		SHR.Apply(sh1b)
		for i := range 3 {
			s.SH1[i*3] = cmn.EncodeSplatSH(float64(sh1r[i]))
			s.SH1[i*3+1] = cmn.EncodeSplatSH(float64(sh1g[i]))
			s.SH1[i*3+2] = cmn.EncodeSplatSH(float64(sh1b[i]))
		}
	}
}

func (s *SplatData) ToString() string {
	return fmt.Sprintf("%v, %v, %v; %v, %v, %v; %v, %v, %v, %v; %v, %v, %v, %v",
		s.PositionX, s.PositionY, s.PositionZ, s.ScaleX, s.ScaleY, s.ScaleZ, s.ColorR, s.ColorG, s.ColorB, s.ColorA, s.RotationW, s.RotationX, s.RotationY, s.RotationZ)
}

func Sort(rows []*SplatData) {
	// PLY没有压缩，忽略排序
	if IsOutputSplat() {
		SortSplat(rows) // 仅编码，按原作排序
	} else if IsOutputSpx() || IsOutputSpz() {
		SortMorton(rows) // 莫顿码排序，提高压缩率
	}
}

func SortSplat(rows []*SplatData) {
	// from https://github.com/antimatter15/splat/blob/main/convert.py
	sort.Slice(rows, func(i, j int) bool {
		return math.Exp(float64(cmn.EncodeSplatScale(rows[i].ScaleX)+cmn.EncodeSplatScale(rows[i].ScaleY)+cmn.EncodeSplatScale(rows[i].ScaleZ)))/(1.0+math.Exp(float64(rows[i].ColorA))) <
			math.Exp(float64(cmn.EncodeSplatScale(rows[j].ScaleX)+cmn.EncodeSplatScale(rows[j].ScaleY)+cmn.EncodeSplatScale(rows[j].ScaleZ)))/(1.0+math.Exp(float64(rows[i].ColorA)))
	})
}

func SortMorton(rows []*SplatData) {
	mm := ComputeXyzMinMax(rows)
	sort.Slice(rows, func(i, j int) bool {
		return EncodeMorton3(rows[i].PositionX, rows[i].PositionY, rows[i].PositionZ, mm) < EncodeMorton3(rows[j].PositionX, rows[j].PositionY, rows[j].PositionZ, mm)
	})
}

func SortBlockDatas4Compress(rows []*SplatData) {
	// sort.Slice(rows, func(i, j int) bool {
	// 	if rows[i].PositionX < rows[j].PositionX {
	// 		return true
	// 	}
	// 	if rows[i].PositionX > rows[j].PositionX {
	// 		return false
	// 	}
	// 	if rows[i].PositionY < rows[j].PositionY {
	// 		return true
	// 	}
	// 	if rows[i].PositionY > rows[j].PositionY {
	// 		return false
	// 	}
	// 	return rows[i].PositionZ < rows[j].PositionZ
	// })
}

func getRotateArgs() (bool, float32, float32, float32) {
	has := Args.HasArgIgnorecase("-rx", "--rotateX", "-ry", "--rotateY", "-rz", "--rotateZ")
	var rx, ry, rz float32
	if has {
		rx = cmn.StringToFloat32(Args.GetArgIgnorecase("-rx", "--rotateX"), 0)
		ry = cmn.StringToFloat32(Args.GetArgIgnorecase("-ry", "--rotateY"), 0)
		rz = cmn.StringToFloat32(Args.GetArgIgnorecase("-rz", "--rotateZ"), 0)
	}
	return has, rx, ry, rz
}

func getScaleArgs() (bool, float32) {
	has := Args.HasArgIgnorecase("-s", "--scale")
	var scale float32 = 1.0
	if has {
		scale = min(max(cmn.StringToFloat32(Args.GetArgIgnorecase("-s", "--scale"), 1.0), 0.001), 1000.0)
	}
	return has, scale
}

func getTranslateArgs() (bool, float32, float32, float32) {
	has := Args.HasArgIgnorecase("-tx", "--translateX", "-ty", "--translateY", "-tz", "--translateZ")
	var tx, ty, tz float32
	if has {
		tx = cmn.StringToFloat32(Args.GetArgIgnorecase("-tx", "--translateX"), 0)
		ty = cmn.StringToFloat32(Args.GetArgIgnorecase("-ty", "--translateY"), 0)
		tz = cmn.StringToFloat32(Args.GetArgIgnorecase("-tz", "--translateZ"), 0)
	}
	return has, tx, ty, tz
}

// ------------- Quaternion --------------
func NewQuaternion(x, y, z, w float64) *Quaternion {
	return &Quaternion{x, y, z, w}

}

// Quaternion :
type Quaternion struct {
	X float64
	Y float64
	Z float64
	W float64
}

func (q *Quaternion) SetFromAxisAngle(axis *Vector3, angle float64) *Quaternion {
	// from http://www.euclideanspace.com/maths/geometry/rotations/conversions/angleToQuaternion/index.htm

	// assumes axis is normalized
	halfAngle := angle / 2
	s := math.Sin(halfAngle)

	q.X = axis.X * s
	q.Y = axis.Y * s
	q.Z = axis.Z * s
	q.W = math.Cos(halfAngle)

	return q
}

func (q *Quaternion) Length() float64 {
	return math.Sqrt(q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W)
}

func (q *Quaternion) Normalize() *Quaternion {
	l := q.Length()

	if l == 0 {
		q.X = 0
		q.Y = 0
		q.Z = 0
		q.W = 1
	} else {
		l = 1 / l

		q.X = q.X * l
		q.Y = q.Y * l
		q.Z = q.Z * l
		q.W = q.W * l
	}

	return q
}

func (q *Quaternion) Multiply(q1 *Quaternion) *Quaternion {
	return q.MultiplyQuaternions(q, q1)
}

func (q *Quaternion) Premultiply(q1 *Quaternion) *Quaternion {
	return q.MultiplyQuaternions(q1, q)
}

func (q *Quaternion) MultiplyQuaternions(a, b *Quaternion) *Quaternion {
	// from http://www.euclideanspace.com/maths/algebra/realNormedAlgebra/quaternions/code/index.htm
	qax, qay, qaz, qaw := a.X, a.Y, a.Z, a.W
	qbx, qby, qbz, qbw := b.X, b.Y, b.Z, b.W

	q.X = qax*qbw + qaw*qbx + qay*qbz - qaz*qby
	q.Y = qay*qbw + qaw*qby + qaz*qbx - qax*qbz
	q.Z = qaz*qbw + qaw*qbz + qax*qby - qay*qbx
	q.W = qaw*qbw - qax*qbx - qay*qby - qaz*qbz

	return q
}

// ------------- Vector3 --------------
func NewVector3(x, y, z float64) *Vector3 {
	return &Vector3{x, y, z}
}

// Vector3 :
type Vector3 struct {
	X float64
	Y float64
	Z float64
}

func (v *Vector3) ApplyQuaternion(q *Quaternion) *Vector3 {
	x, y, z := v.X, v.Y, v.Z
	qx, qy, qz, qw := q.X, q.Y, q.Z, q.W

	// calculate quat * vector

	ix := qw*x + qy*z - qz*y
	iy := qw*y + qz*x - qx*z
	iz := qw*z + qx*y - qy*x
	iw := -qx*x - qy*y - qz*z

	// calculate result * inverse quat

	v.X = ix*qw + iw*-qx + iy*-qz - iz*-qy
	v.Y = iy*qw + iw*-qy + iz*-qx - ix*-qz
	v.Z = iz*qw + iw*-qz + ix*-qy - iy*-qx

	return v
}

func CompressionInfo(filePath string, num int, inFileSize ...int64) string {
	if cmn.Endwiths(filePath, ".ply", true) && !cmn.Endwiths(filePath, ".compressed.ply", true) {
		return fmt.Sprintf("splat count: %v", num)
	}

	fileSize := cmn.GetFileSize(filePath)
	if cmn.FileName(filePath) == "meta.json" && len(inFileSize) > 0 {
		fileSize = inFileSize[0] // sog 索引文件时，使用参数传入的文件大小值
	}

	plySize := 1500 + num*248
	compressionRatio := float64(plySize) / float64(fileSize)
	sizeReduction := (1 - float64(fileSize)/float64(plySize)) * 100
	fileSizeM := float64(fileSize) / 1024.0 / 1024.0

	shDegree := GetArgShDegree()
	return fmt.Sprintf("splat count: %v, %.1fM, %.2fx compression with sh%v (%.2f%% smaller than 3dgs ply)", num, fileSizeM, compressionRatio, shDegree, sizeReduction)
}
