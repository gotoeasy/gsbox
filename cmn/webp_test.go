package cmn

import (
	"log"
	"testing"
	"time"
)

func Test_webp(t *testing.T) {
	fileWebp := "e:\\test.webp"
	webpBytes, err := ReadFileBytes(fileWebp)
	if err != nil {
		t.Fatal(err)
	}

	rgbas, w, h, err := DecompressWebp(webpBytes)
	if err != nil {
		t.Fatal(err)
	}

	// compress by gen2brain/webp
	startTime := time.Now()
	rsWebpBytes, err := CompressWebpByWidthHeight(rgbas, w, h, 90)
	if err != nil {
		t.Fatal(err)
	}

	for range 9 {
		rsWebpBytes, err = CompressWebpByWidthHeight(rgbas, w, h, 90)
		if err != nil {
			t.Fatal(err)
		}
	}
	log.Println("[Info] processing time:", GetTimeInfo(time.Since(startTime).Milliseconds()))

	// print log when dynamic
	PrintLibwebpInfo(true)

	// write to file
	err = WriteFileBytes(fileWebp+".webp", rsWebpBytes)
	if err != nil {
		t.Fatal(err)
	}
}
