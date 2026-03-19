package gsplat

import (
	"bufio"
	"fmt"
	"gsbox/cmn"
	"log"
	"os"
	"strings"
)

// 当前仅用于顶点坐标转换
func Obj2Obj(srcPath, dstPath string) {
	srcFile, err := os.Open(srcPath)
	cmn.ExitOnError(err)
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	cmn.ExitOnError(err)
	defer dstFile.Close()

	scanner := bufio.NewScanner(srcFile)
	writer := bufio.NewWriter(dstFile)

	hasRotate, degreeX, degreeY, degreeZ := getRotateArgs()
	hasScale, scale := getScaleArgs()
	hasTranslate, tx, ty, tz := getTranslateArgs()
	order := cmn.ToLower(Args.GetArgIgnorecase("-to", "--transform-order"))

	if hasRotate {
		log.Println("[Info] (transform) rotate in XYZ order.", "degreeX:", degreeX, ", degreeY:", degreeY, ", degreeZ:", degreeZ)
	}
	if hasScale {
		log.Println("[Info] (transform) scaling factor:", scale)
	}
	if hasTranslate {
		log.Println("[Info] (transform) make translate.", "translateX:", tx, ", translateY:", ty, ", translateZ:", tz)
	}

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if cmn.Startwiths(line, "v ") && hasRotate {
			parts := strings.Fields(line)
			x, y, z := cmn.StringToFloat32(parts[1]), cmn.StringToFloat32(parts[2]), cmn.StringToFloat32(parts[3])
			x, y, z = transformVertex(x, y, z, degreeX, degreeY, degreeZ, scale, tx, ty, tz, hasRotate, hasScale, hasTranslate, order)
			newLine := fmt.Sprintf("v %.6f %.6f %.6f", x, y, z)
			_, err = writer.WriteString(newLine + "\n")
			cmn.ExitOnError(err)
		} else {
			_, err = writer.WriteString(line + "\n")
			cmn.ExitOnError(err)
		}
	}

	cmn.ExitOnError(scanner.Err())
	cmn.ExitOnError(writer.Flush())
}

func transformVertex(x, y, z, degreeX, degreeY, degreeZ, scale, tx, ty, tz float32, hasRotate, hasScale, hasTranslate bool, order string) (newX, newY, newZ float32) {
	newX, newY, newZ = x, y, z
	if !hasRotate && !hasScale && !hasTranslate {
		return
	}

	switch order {
	case "rts":
		newX, newY, newZ = rotateVertex(newX, newY, newZ, degreeX, degreeY, degreeZ, hasRotate)
		newX, newY, newZ = translateVertex(newX, newY, newZ, tx, ty, tz, hasTranslate)
		newX, newY, newZ = scaleVertex(newX, newY, newZ, scale, hasScale)
	case "srt":
		newX, newY, newZ = scaleVertex(newX, newY, newZ, scale, hasScale)
		newX, newY, newZ = rotateVertex(newX, newY, newZ, degreeX, degreeY, degreeZ, hasRotate)
		newX, newY, newZ = translateVertex(newX, newY, newZ, tx, ty, tz, hasTranslate)
	case "str":
		newX, newY, newZ = scaleVertex(newX, newY, newZ, scale, hasScale)
		newX, newY, newZ = translateVertex(newX, newY, newZ, tx, ty, tz, hasTranslate)
		newX, newY, newZ = rotateVertex(newX, newY, newZ, degreeX, degreeY, degreeZ, hasRotate)
	case "trs":
		newX, newY, newZ = translateVertex(newX, newY, newZ, tx, ty, tz, hasTranslate)
		newX, newY, newZ = rotateVertex(newX, newY, newZ, degreeX, degreeY, degreeZ, hasRotate)
		newX, newY, newZ = scaleVertex(newX, newY, newZ, scale, hasScale)
	case "tsr":
		newX, newY, newZ = translateVertex(newX, newY, newZ, tx, ty, tz, hasTranslate)
		newX, newY, newZ = scaleVertex(newX, newY, newZ, scale, hasScale)
		newX, newY, newZ = rotateVertex(newX, newY, newZ, degreeX, degreeY, degreeZ, hasRotate)
	default:
		newX, newY, newZ = rotateVertex(newX, newY, newZ, degreeX, degreeY, degreeZ, hasRotate)
		newX, newY, newZ = scaleVertex(newX, newY, newZ, scale, hasScale)
		newX, newY, newZ = translateVertex(newX, newY, newZ, tx, ty, tz, hasTranslate)
	}
	return
}

func rotateVertex(x, y, z, degreeX, degreeY, degreeZ float32, hasRotate bool) (newX, newY, newZ float32) {
	newX, newY, newZ = x, y, z
	if !hasRotate {
		return
	}

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
	point := NewVector3(float64(x), float64(y), float64(z))
	point.ApplyQuaternion(q)
	newX, newY, newZ = cmn.ClipFloat32(point.X), cmn.ClipFloat32(point.Y), cmn.ClipFloat32(point.Z)
	return
}

func translateVertex(x, y, z, tx, ty, tz float32, hasTranslate bool) (newX, newY, newZ float32) {
	newX, newY, newZ = x, y, z
	if !hasTranslate {
		return
	}
	newX, newY, newZ = x+tx, y+ty, z+tz
	return
}

func scaleVertex(x, y, z, scale float32, hasScale bool) (newX, newY, newZ float32) {
	newX, newY, newZ = x, y, z
	if !hasScale {
		return
	}
	newX, newY, newZ = x*scale, y*scale, z*scale
	return
}
