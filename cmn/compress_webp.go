package cmn

import (
	"bytes"
	"errors"
	"image"
	"log"
	"math"

	"github.com/HugoSmits86/nativewebp"
	"github.com/gen2brain/webp"
)

func CompressWebpByWidthHeight(bts []byte, width int, height int) ([]byte, error) {
	datas := bts
	if len(bts)/4 < width*height {
		datas = append(datas, bytes.Repeat([]uint8{0}, width*height*4-len(bts))...)
	}

	var buf bytes.Buffer
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	copy(img.Pix, datas)
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

func CompressWebp(bts []byte) ([]byte, error) {
	width, height := ComputeWidthHeight(len(bts))
	return CompressWebpByWidthHeight(bts, width, height)
}

func DecompressWebp(webpBytes []byte) (rgbas []byte, width int, height int, err error) {
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

func ComputeWidthHeight(length int) (width int, height int) {
	w := math.Ceil(math.Sqrt(float64(length))/4.0) * 4.0
	h := math.Ceil(float64(length)/w/4.0) * 4.0
	return int(w), int(h)
}

func PrintLibwebpInfo(hasWebp bool) {
	if hasWebp && webp.Dynamic() == nil {
		log.Println("[Info] using libwebp for webp compression - perfect setup")
	}
}
