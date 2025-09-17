package cmn

import (
	"fmt"
	"log"
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

func Test_xyz(t *testing.T) {
	value := float32(1000.12345678987654321)
	encodeValue := EncodeXyz(value)
	decodeValue := DecodeXyz(encodeValue)
	log.Println(value, encodeValue, decodeValue)

}
