package utils

import (
	"bytes"
	"compress/gzip"
)

// Gzip 压缩
func GzipEncode(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)

	if _, err := gzw.Write(input); err != nil {
		return nil, err
	}
	if err := gzw.Flush(); err != nil {
		return nil, err
	}
	if err := gzw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Gzip 解压
func GzipDecode(input []byte) ([]byte, error) {
	gzr, err := gzip.NewReader(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	var output = new(bytes.Buffer)
	if _, err := output.ReadFrom(gzr); err != nil {
		return nil, err
	}
	if err := gzr.Close(); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}
