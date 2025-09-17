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
	bs, _, _, err := DecompressWebpMore(webpBytes)
	return bs, err
}

func DecompressWebpMore(webpBytes []byte) ([]byte, int, int, error) {
	reader := bytes.NewReader(webpBytes)
	decodedImg, err := nativewebp.Decode(reader)
	if err != nil {
		return nil, 0, 0, err
	}

	nrgba, ok := decodedImg.(*image.NRGBA)
	if !ok {
		return nil, 0, 0, errors.New("decoded image is not *image.NRGBA")
	}
	return nrgba.Pix, nrgba.Bounds().Size().X, nrgba.Bounds().Size().Y, nil
}

func computeWidthHeight(length int) (int, int) {
	pixels := math.Ceil(float64(length) / 4.0)
	width := math.Ceil(math.Sqrt(pixels))
	height := (pixels + width - 1) / width
	return int(width), int(height)
}
