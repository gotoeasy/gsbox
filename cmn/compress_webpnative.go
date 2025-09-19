//go:build !cgo

package cmn

import (
	"bytes"
	"image"

	"github.com/HugoSmits86/nativewebp"
)

func (c *nativeCompressor) Compress(data []byte, widthHeight ...int) ([]byte, error) {
	var width, height int
	if len(widthHeight) == 2 {
		width = widthHeight[0]
		height = widthHeight[1]
	} else {
		width, height = ComputeWidthHeight(len(data))
	}
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	copy(img.Pix, data)

	var buf bytes.Buffer
	err := nativewebp.Encode(&buf, img, nil)
	return buf.Bytes(), err
}

type nativeCompressor struct{}

func init() {
	webpComp = &nativeCompressor{}
}
