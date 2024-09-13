package gsplat

import (
	"bufio"
	"gsbox/cmn"
	"os"
)

func WriteSplat(splatFile string, rows []*SplatData, headers ...string) {
	file, err := os.Create(splatFile)
	cmn.ExitOnError(err)
	defer file.Close()

	writer := bufio.NewWriter(file)

	// 自定义头
	if len(headers) > 0 && headers[0] != "" {
		writer.WriteString(headers[0])
	}

	// 内容
	for i := 0; i < len(rows); i++ {
		_, err = writer.Write(rows[i].ToBytes())
		cmn.ExitOnError(err)
	}
	err = writer.Flush()
	cmn.ExitOnError(err)
}
