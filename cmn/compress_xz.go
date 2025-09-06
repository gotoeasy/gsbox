package cmn

import (
	"bytes"
	"io"

	"github.com/ulikunitz/xz"
)

func CompressXZ(bts []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := xz.NewWriter(&buf)
	if err != nil {
		return nil, err
	}

	_, err = w.Write(bts)
	if err != nil {
		return nil, err
	}
	w.Close()

	return buf.Bytes(), nil
}

func DecompressXZ(xzBytes []byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write(xzBytes)

	r, err := xz.NewReader(&buf)
	if err != nil {
		return nil, err
	}

	decompressedBuf := &bytes.Buffer{}
	_, err = io.Copy(decompressedBuf, r)
	if err != nil {
		return nil, err
	}

	return decompressedBuf.Bytes(), nil
}
