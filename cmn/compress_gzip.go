package cmn

import (
	"bytes"
	"compress/gzip"
	"io"
)

// 用gzip压缩字节数组
func GzipBytes(bts []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write(bts); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 解压gzip字节数组
func UnGzipBytes(gzipBytes []byte) ([]byte, error) {
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
