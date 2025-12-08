package gsplat

type SogHeader struct {
	Version     int
	Count       int
	PaletteSize int
	ShDegree    uint8

	Palettes []uint8
}
