package cmn

import (
	"bytes"
	"errors"
	"image"
	"math"

	"github.com/HugoSmits86/nativewebp"
)

func CompressWebp(bts []byte) ([]byte, error) {

	length := len(bts)
	width, height := computeWidthHeight(length)
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	copy(img.Pix, bts)

	var buf bytes.Buffer
	err := nativewebp.Encode(&buf, img, nil)
	return buf.Bytes(), err
}

func DecompressWebp(webpBytes []byte) ([]byte, error) {
	reader := bytes.NewReader(webpBytes)
	decodedImg, err := nativewebp.Decode(reader)
	if err != nil {
		return nil, err
	}

	nrgba, ok := decodedImg.(*image.NRGBA)
	if !ok {
		return nil, errors.New("decoded image is not *image.NRGBA")
	}
	return nrgba.Pix, nil
}

func computeWidthHeight(length int) (int, int) {
	pixels := math.Ceil(float64(length) / 4.0)
	width := math.Ceil(math.Sqrt(pixels))
	height := (pixels + width - 1) / width
	return int(width), int(height)
}
