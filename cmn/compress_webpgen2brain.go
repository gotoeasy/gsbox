package cmn

import (
	"bytes"
	"image"

	"github.com/gen2brain/webp"
)

func (c *gen2brainCompressor) Compress(data []byte, width int, height int) ([]byte, error) {
	var buf bytes.Buffer
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	copy(img.Pix, data)
	options := webp.Options{
		Quality:  90,    // 质量，默认75，最大100
		Lossless: true,  // 无损
		Method:   6,     // 越大压缩效果越好速度最慢，默认4，最大6
		Exact:    false, // 透明时不需要保持RGB值
	}
	err := webp.Encode(&buf, img, options)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type gen2brainCompressor struct{}

func init() {
	webpComp = &gen2brainCompressor{}
}
