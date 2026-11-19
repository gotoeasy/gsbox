package gsplat

import (
	"errors"
	"gsbox/cmn"
)

func WriteLodMetaSog(output string, splatTiles *SplatTiles, lodMeta *LodMeta) {

	ver := Args.GetArgIgnorecase("-ov", "--output-version")

	switch ver {
	case "", "3":
		writeLodMetaSogV3(output, splatTiles, lodMeta)
	case "2":
		writeLodMetaSogV2(output, splatTiles, lodMeta)
	default:
		cmn.ExitOnError(errors.New("unsupported sog version: " + ver))
	}
}
