package cmn

import (
	"fmt"
	"testing"
)

func Test_cmn(t *testing.T) {
	value := float32(12.3456789)

	encodedBytes := EncodeFloat32ToBytes3(value)
	fmt.Printf("Encoded bytes: %v\n", encodedBytes)

	decodedValue := DecodeBytes3ToFloat32(encodedBytes)
	fmt.Printf("Decoded value: %f\n", decodedValue)

	b := EncodeFloat32ToByte(value)
	fmt.Printf("Encoded byte: %v\n", b)
	d := DecodeByteToFloat32(b)
	fmt.Printf("Decoded value: %f\n", d)

	fmt.Println("-------------------")
	value = float32(-12.3456789)

	encodedBytes = EncodeFloat32ToBytes3(value)
	fmt.Printf("Encoded bytes: %v\n", encodedBytes)

	decodedValue = DecodeBytes3ToFloat32(encodedBytes)
	fmt.Printf("Decoded value: %f\n", decodedValue)
}

// func Test_zstd(t *testing.T) {
// 	var bts []byte
// 	for i := range 200 {
// 		bts = append(bts, uint8(i))
// 	}

// 	bs, err := ZstdCompress(bts)
// 	log.Println("bs length", len(bs), err)

// 	bs2, err := ZstdDecompress(bs)
// 	log.Println("bs2 length", len(bs2), bs2, err)

// }
