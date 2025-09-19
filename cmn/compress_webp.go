package cmn

import (
	"bytes"
	"errors"
	"image"
	"math"

	"github.com/HugoSmits86/nativewebp"
)

func CompressWebp(bts []byte, widthHeight ...int) ([]byte, error) {
	return webpComp.Compress(bts, widthHeight...)
}

func DecompressWebp(webpBytes []byte) ([]byte, int, int, error) {
	reader := bytes.NewReader(webpBytes)
	img, err := nativewebp.Decode(reader)
	if err != nil {
		return nil, 0, 0, err
	}

	nrgba, ok := img.(*image.NRGBA)
	if !ok {
		return nil, 0, 0, errors.New("decoded image is not *image.NRGBA")
	}
	return nrgba.Pix, img.Bounds().Size().X, img.Bounds().Size().Y, nil
}

func ComputeWidthHeight(length int) (int, int) {
	pixels := math.Ceil(float64(length) / 4.0)
	width := math.Ceil(math.Sqrt(pixels))
	height := (pixels + width - 1) / width
	return int(width), int(height)
}

var webpComp webpCompressor

type webpCompressor interface {
	Compress(data []byte, widthHeight ...int) ([]byte, error)
}
