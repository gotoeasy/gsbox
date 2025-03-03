package gsplat

import (
	"bufio"
	"gsbox/cmn"
	"os"
)

func WriteSplat20(splat20File string, rows []*SplatData) {
	file, err := os.Create(splat20File)
	cmn.ExitOnError(err)
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i := 0; i < len(rows); i++ {
		_, err = writer.Write(rows[i].ToBytesSplat20())
		cmn.ExitOnError(err)
	}
	err = writer.Flush()
	cmn.ExitOnError(err)
}
