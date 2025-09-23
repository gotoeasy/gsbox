package cmn

import (
	"bytes"
	"errors"
	"image"
	"math"

	"github.com/HugoSmits86/nativewebp"
)

func CompressWebp(bts []byte) ([]byte, error) {
	width, height := ComputeWidthHeight(len(bts))
	if len(bts) == width*height {
		return webpComp.Compress(bts, width, height)
	}

	// 数据少的没必要webp编码压缩，只考虑最终图片至少32*32以上大小的场景
	datas := bts[:]
	fullRows := len(bts) / width
	partRowLen := len(bts) % width
	dataRows := fullRows
	if partRowLen > 0 {
		dataRows++
	}
	if partRowLen > 0 {
		// 不足一行的，取上行同列填充
		start := fullRows*width - width + partRowLen
		addPartRow := bts[start : start+width-partRowLen]
		datas = append(datas, addPartRow...)
	}

	addFullRowCnt := height - dataRows // max 3
	for i := range addFullRowCnt {
		// 空行的，倒序逐行取整行填充
		srcRow := dataRows - 1 - i
		if partRowLen > 0 {
			srcRow--
		}
		addFullRow := bts[(srcRow-1)*width : srcRow*width]
		datas = append(datas, addFullRow...)
	}

	return webpComp.Compress(datas, width, height)
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

var webpComp webpCompressor

type webpCompressor interface {
	Compress(data []byte, widthHeight ...int) ([]byte, error)
}
