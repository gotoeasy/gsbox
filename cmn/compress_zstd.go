package cmn

// zstd 效果 还不如 gzip 的样子，暂不支持

// import (
// 	"github.com/klauspost/compress/zstd"
// )

// // 用zstd压缩字节数组
// func ZstdCompress(bts []byte) ([]byte, error) {
// 	zstdEncoder, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBetterCompression))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer zstdEncoder.Close()

// 	return zstdEncoder.EncodeAll(bts, make([]byte, 0, len(bts))), nil
// }

// // 解压zstd字节数组
// func ZstdDecompress(zstdBytes []byte) ([]byte, error) {
// 	zstdDecoder, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer zstdDecoder.Close()

// 	return zstdDecoder.DecodeAll(zstdBytes, nil)
// }
