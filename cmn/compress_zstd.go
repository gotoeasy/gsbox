package cmn

import (
	"github.com/klauspost/compress/zstd"
)

func CompressZstd(bts []byte) ([]byte, error) {
	zstdEncoder, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBetterCompression))
	if err != nil {
		return nil, err
	}
	defer zstdEncoder.Close()

	return zstdEncoder.EncodeAll(bts, make([]byte, 0, len(bts))), nil
}

func DecompressZstd(zstdBytes []byte) ([]byte, error) {
	zstdDecoder, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))
	if err != nil {
		return nil, err
	}
	defer zstdDecoder.Close()

	return zstdDecoder.DecodeAll(zstdBytes, nil)
}
