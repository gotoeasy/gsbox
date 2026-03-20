package gsplat

import (
	"encoding/json"
	"gsbox/cmn"
	"os"
)

func ReadGlbJson(glbFile string) string {
	file, err := os.Open(glbFile)
	cmn.ExitOnError(err)
	defer file.Close()

	bs := make([]byte, 20)
	_, err = file.Read(bs)
	cmn.ExitOnError(err)

	jsonLength := int(cmn.BytesToUint32(bs[12:16]))

	bs = make([]byte, jsonLength)
	_, err = file.Read(bs)
	cmn.ExitOnError(err)

	strJson := cmn.BytesToString(bs)

	var data any
	cmn.ExitOnError(json.Unmarshal([]byte(strJson), &data))

	bts, err := json.MarshalIndent(data, "", "  ")
	cmn.ExitOnError(err)
	return cmn.BytesToString(bts)
}
