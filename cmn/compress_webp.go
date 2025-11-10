package cmn

import (
	"bytes"
	"errors"
	"image"
	"math"

	"github.com/HugoSmits86/nativewebp"
)

func CompressWebpByWidthHeight(bts []byte, width int, height int) ([]byte, error) {
	datas := bts
	if len(bts) != width*height {
		datas = bts[:]
		dataPixCnt := len(bts) / 4
		fullRows := dataPixCnt / width                                     // 完整数据行的数量
		partRowDataPix := dataPixCnt % width                               // 部分数据行中数据的像素数量
		dataRows := int(math.Ceil((float64(dataPixCnt) / float64(width)))) // 总数据行数
		if partRowDataPix > 0 {
			// 部分数据行自动填充
			if dataRows > 1 {
				// 多行，取上行同列填充
				start := fullRows*width - width + partRowDataPix
				addPartRow := bts[start : start+width-partRowDataPix]
				datas = append(datas, addPartRow...)
			} else {
				// 仅1行，用最后数据像素填充
				pix := bts[dataPixCnt*4-4 : dataPixCnt*4]
				for i := dataPixCnt; i < width; i++ {
					datas = append(datas, pix...)
				}
			}
		}

		// 参数指定行比数据行大，用最后行补足
		if height > dataRows {
			addRowCnt := height - dataRows
			row := datas[fullRows*width*4 : (fullRows+1)*width*4]
			for range addRowCnt {
				datas = append(datas, row...)
			}
		}
	}

	return webpComp.Compress(datas, width, height)
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

var webpComp webpCompressor

type webpCompressor interface {
	Compress(data []byte, width int, height int) ([]byte, error)
}
