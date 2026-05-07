package cmn

import (
	"github.com/klauspost/compress/zstd"
)

// 全局复用的编码器和解码器（线程安全，可并发调用）
var (
	encoderZstd, _ = zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	decoderZstd, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))
)

func CompressZstd(bts []byte) []byte {
	return encoderZstd.EncodeAll(bts, make([]byte, 0, len(bts)))
}

func DecompressZstd(zstdBytes []byte) ([]byte, error) {
	return decoderZstd.DecodeAll(zstdBytes, nil)
}
