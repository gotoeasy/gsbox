//go:build cgo

package cmn

import (
	"image"

	"github.com/chai2010/webp"
)

func (c *cgoCompressor) Compress(data []byte, widthHeight ...int) ([]byte, error) {
	var width, height int
	if len(widthHeight) == 2 {
		width = widthHeight[0]
		height = widthHeight[1]
	} else {
		width, height = ComputeWidthHeight(len(data))
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	copy(img.Pix, data)
	return webp.EncodeLosslessRGBA(img)
}

type cgoCompressor struct{}

func init() {
	webpComp = &cgoCompressor{}
}
