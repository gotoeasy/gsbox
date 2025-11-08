package gsplat

import (
	"gsbox/cmn"
	"math"
)

type V3MinMax struct {
	MinX    float32
	MinY    float32
	MinZ    float32
	MaxX    float32
	MaxY    float32
	MaxZ    float32
	LenX    float32
	LenY    float32
	LenZ    float32
	CenterX float32
	CenterY float32
	CenterZ float32
}

func ComputeXyzMinMax(datas []*SplatData) *V3MinMax {
	xyzRange := &V3MinMax{
		MinX: math.MaxFloat32,
		MinY: math.MaxFloat32,
		MinZ: math.MaxFloat32,
		MaxX: -math.MaxFloat32,
		MaxY: -math.MaxFloat32,
		MaxZ: -math.MaxFloat32,
	}

	for i := range datas {
		x := datas[i].PositionX
		y := datas[i].PositionY
		z := datas[i].PositionZ
		xyzRange.MinX = min(xyzRange.MinX, x)
		xyzRange.MinY = min(xyzRange.MinY, y)
		xyzRange.MinZ = min(xyzRange.MinZ, z)
		xyzRange.MaxX = max(xyzRange.MaxX, x)
		xyzRange.MaxY = max(xyzRange.MaxY, y)
		xyzRange.MaxZ = max(xyzRange.MaxZ, z)
	}

	xyzRange.LenX = xyzRange.MaxX - xyzRange.MinX
	xyzRange.LenY = xyzRange.MaxY - xyzRange.MinY
	xyzRange.LenZ = xyzRange.MaxZ - xyzRange.MinZ
	xyzRange.CenterX = (xyzRange.MaxX + xyzRange.MinX) / 2.0
	xyzRange.CenterY = (xyzRange.MaxY + xyzRange.MinY) / 2.0
	xyzRange.CenterZ = (xyzRange.MaxZ + xyzRange.MinZ) / 2.0

	return xyzRange
}

func ComputeXyzLogMinMax(datas []*SplatData) *V3MinMax {
	xyzRange := &V3MinMax{
		MinX: math.MaxFloat32,
		MinY: math.MaxFloat32,
		MinZ: math.MaxFloat32,
		MaxX: -math.MaxFloat32,
		MaxY: -math.MaxFloat32,
		MaxZ: -math.MaxFloat32,
	}

	for i := range datas {
		x := cmn.SogEncodeLog(datas[i].PositionX)
		y := cmn.SogEncodeLog(datas[i].PositionY)
		z := cmn.SogEncodeLog(datas[i].PositionZ)
		xyzRange.MinX = min(xyzRange.MinX, x)
		xyzRange.MinY = min(xyzRange.MinY, y)
		xyzRange.MinZ = min(xyzRange.MinZ, z)
		xyzRange.MaxX = max(xyzRange.MaxX, x)
		xyzRange.MaxY = max(xyzRange.MaxY, y)
		xyzRange.MaxZ = max(xyzRange.MaxZ, z)
	}

	xyzRange.LenX = xyzRange.MaxX - xyzRange.MinX
	xyzRange.LenY = xyzRange.MaxY - xyzRange.MinY
	xyzRange.LenZ = xyzRange.MaxZ - xyzRange.MinZ
	xyzRange.CenterX = (xyzRange.MaxX + xyzRange.MinX) / 2.0
	xyzRange.CenterY = (xyzRange.MaxY + xyzRange.MinY) / 2.0
	xyzRange.CenterZ = (xyzRange.MaxZ + xyzRange.MinZ) / 2.0

	return xyzRange
}

// https://fgiesen.wordpress.com/2009/12/13/decoding-morton-codes/
func EncodeMorton3(x, y, z float32, mm *V3MinMax) uint32 {
	ix := min(1023, uint32(math.Floor(float64(1024.0*(x-mm.MinX)/mm.LenX))))
	iy := min(1023, uint32(math.Floor(float64(1024.0*(y-mm.MinY)/mm.LenY))))
	iz := min(1023, uint32(math.Floor(float64(1024.0*(z-mm.MinZ)/mm.LenZ))))

	return (Part1By2(iz) << 2) + (Part1By2(iy) << 1) + Part1By2(ix)
}

func Part1By1(x uint32) uint32 {
	x &= 0x0000ffff                 // x = ---- ---- ---- ---- fedc ba98 7654 3210
	x = (x ^ (x << 8)) & 0x00ff00ff // x = ---- ---- fedc ba98 ---- ---- 7654 3210
	x = (x ^ (x << 4)) & 0x0f0f0f0f // x = ---- fedc ---- ba98 ---- 7654 ---- 3210
	x = (x ^ (x << 2)) & 0x33333333 // x = --fe --dc --ba --98 --76 --54 --32 --10
	x = (x ^ (x << 1)) & 0x55555555 // x = -f-e -d-c -b-a -9-8 -7-6 -5-4 -3-2 -1-0
	return x
}

func Part1By2(x uint32) uint32 {
	x &= 0x000003FF
	x = (x ^ (x << 16)) & 0xFF0000FF
	x = (x ^ (x << 8)) & 0x0300F00F
	x = (x ^ (x << 4)) & 0x030C30C3
	x = (x ^ (x << 2)) & 0x09249249
	return x
}
