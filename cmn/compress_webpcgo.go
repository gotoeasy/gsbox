//go:build cgo

package cmn

import (
	"image"

	"github.com/chai2010/webp"
)

func (c *cgoCompressor) Compress(data []byte, width int, height int) (webps []byte, err error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	copy(img.Pix, data)
	return webp.EncodeLosslessRGBA(img)
}

type cgoCompressor struct{}

func init() {
	webpComp = &cgoCompressor{}
}
