package gsplat

import (
	"errors"
	"gsbox/cmn"
)

// 支持输出zip压缩单个sog文件，或压缩前的多个文件
func WriteSog(sogOrJsonFile string, rows []*SplatData) (fileSize int64) {
	ver := Args.GetArgIgnorecase("-ov", "--output-version")

	switch ver {
	case "", "3":
		return WriteSogV3(sogOrJsonFile, rows)
	case "2":
		return WriteSogV2(sogOrJsonFile, rows)
	default:
		cmn.ExitOnError(errors.New("unsupported sog version: " + ver))
	}
	return
}
