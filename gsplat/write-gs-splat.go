package gsplat

import (
	"bufio"
	"gsbox/cmn"
	"os"
)

func WriteSplat(splatFile string, rows []*SplatData) {
	file, err := os.Create(splatFile)
	cmn.ExitOnError(err)
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i := 0; i < len(rows); i++ {
		_, err = writer.Write(rows[i].ToBytes())
		cmn.ExitOnError(err)
	}
	err = writer.Flush()
	cmn.ExitOnError(err)
}
