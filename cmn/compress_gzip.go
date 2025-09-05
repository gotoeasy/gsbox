package cmn

import (
	"bytes"
	"compress/gzip"
	"io"
)

// 用gzip压缩字节数组
func CompressGzip(bts []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	// gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression) // 差异极小
	// if err != nil {
	// 	return nil, err
	// }

	if _, err := gz.Write(bts); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 解压gzip字节数组
func DecompressGzip(gzipBytes []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(gzipBytes))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	unGzipdBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return unGzipdBytes, nil
}
