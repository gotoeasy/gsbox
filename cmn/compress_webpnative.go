package cmn

// import (
// 	"bytes"
// 	"image"

// 	"github.com/HugoSmits86/nativewebp"
// )

// func (c *nativeCompressor) Compress(data []byte, width int, height int) ([]byte, error) {
// 	img := image.NewNRGBA(image.Rect(0, 0, width, height))
// 	copy(img.Pix, data)

// 	var buf bytes.Buffer
// 	err := nativewebp.Encode(&buf, img, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }

// type nativeCompressor struct{}

// func init() {
// 	webpComp = &nativeCompressor{}
// }
