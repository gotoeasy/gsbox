package gsplat

import (
	"fmt"
	"gsbox/cmn"
	"log"
)

const NewestSpxVersion = int(3)

func OutputSpxVersion() int {
	ver := cmn.StringToInt(Args.GetArgIgnorecase("-ov", "--output-version"), 3)
	return max(1, min(ver, NewestSpxVersion))
}

func IsShChanged() bool {
	return Args.HasArgIgnorecase("-rx", "--rotateX", "-ry", "--rotateY", "-rz", "--rotateZ")
}

func IsSpx2Spx() bool {
	return !Args.HasCmd("join") && Args.HasCmd("x2x", "spx2spx")
}

func IsSog2Spx() bool {
	return !Args.HasCmd("join") && Args.HasCmd("g2x", "sog2spx")
}

func IsSpx2Sog() bool {
	return !Args.HasCmd("join") && Args.HasCmd("x2g", "spx2sog")
}

func IsSog2Sog() bool {
	return !Args.HasCmd("join") && Args.HasCmd("g2g", "sog2sog")
}

func IsSpx2Spz() bool {
	return !Args.HasCmd("join") && Args.HasCmd("x2z", "spx2spz")
}

func IsSog2Spz() bool {
	return !Args.HasCmd("join") && Args.HasCmd("g2z", "sog2spz")
}

func IsSog2SpxOrSog() bool {
	return IsSog2Spx() || IsSog2Sog()
}

func IsSpx2SpxOrSog() bool {
	return IsSpx2Sog() || IsSpx2Spx()
}

func IsSpxOrSog2Spx() bool {
	return IsSpx2Spx() || IsSog2Spx()
}

func IsSpxOrSog2Sog() bool {
	return IsSpx2Sog() || IsSog2Sog()
}

func IsOutputSpx() bool {
	if Args.HasCmd("p2x", "ply2spx", "s2x", "splat2spx", "z2x", "spz2spx", "x2x", "spx2spx", "k2x", "ksplat2spx", "g2x", "sog2spx") {
		return true
	}
	if Args.HasCmd("join") && cmn.Endwiths(Args.GetArgIgnorecase("-o", "--output"), ".spx", true) {
		return true
	}
	return false
}

func IsOutputSpz() bool {
	if Args.HasCmd("p2z", "ply2spz", "s2z", "splat2spz", "z2z", "spz2spz", "x2z", "spx2spz", "k2z", "ksplat2spz", "g2z", "sog2spz") {
		return true
	}
	if Args.HasCmd("join") && cmn.Endwiths(Args.GetArgIgnorecase("-o", "--output"), ".spz", true) {
		return true
	}
	return false
}

func IsOutputPly() bool {
	if Args.HasCmd("p2p", "ply2ply", "s2p", "splat2ply", "z2p", "spz2ply", "x2p", "spx2ply", "k2p", "ksplat2ply", "g2p", "sog2ply") {
		return true
	}
	if Args.HasCmd("join") && cmn.Endwiths(Args.GetArgIgnorecase("-o", "--output"), ".ply", true) {
		return true
	}
	return false
}

func IsOutputSplat() bool {
	if Args.HasCmd("p2s", "ply2splat", "s2s", "splat2splat", "z2s", "spz2splat", "x2s", "spx2splat", "k2s", "ksplat2splat", "g2s", "sog2splat") {
		return true
	}
	if Args.HasCmd("join") && cmn.Endwiths(Args.GetArgIgnorecase("-o", "--output"), ".splat", true) {
		return true
	}
	return false
}

func IsOutputSog() bool {
	if Args.HasCmd("p2g", "ply2sog", "s2g", "splat2sog", "z2g", "spz2sog", "x2g", "spx2sog", "k2g", "ksplat2sog", "g2g", "sog2sog") {
		return true
	}
	if Args.HasCmd("join") && cmn.Endwiths(Args.GetArgIgnorecase("-o", "--output"), ".sog", true) {
		return true
	}
	return false
}

func FilterDatas(datas []*SplatData) []*SplatData {

	dataLen := len(datas)
	if dataLen == 0 || !Args.HasArgIgnorecase("-a", "--alpha") {
		return datas
	}

	inputAlpha := cmn.StringToInt(Args.GetArgIgnorecase("-a", "--alpha"), 0)
	if inputAlpha <= 0 {
		return datas // 无可过滤
	}

	if inputAlpha > 255 {
		inputAlpha = 255
	}

	// 最后要兜底检查，不能大到无数据
	alphas := []int{} // 长度256，使最后位的值为0便于比较计算
	for i := 0; i <= 256; i++ {
		alphas = append(alphas, 0)
	}
	for i := range dataLen {
		alphas[datas[i].ColorA] = alphas[datas[i].ColorA] + 1
	}
	for i := 255; i >= 0; i-- {
		alphas[i] = alphas[i] + alphas[i+1]
	}
	var maxAlpha uint8 // 兜底的alpha
	for i := 255; i >= 0; i-- {
		if alphas[i] > 0 {
			maxAlpha = uint8(i)
			break
		}
	}

	var alpha uint8 = uint8(inputAlpha)
	if alpha > maxAlpha {
		alpha = maxAlpha
	}

	// 过滤数据
	rs := []*SplatData{}
	for i := range dataLen {
		if datas[i].ColorA >= alpha {
			rs = append(rs, datas[i])
		}
	}

	// 日志
	if dataLen-len(rs) > 0 {
		removed := cmn.IntToString(dataLen-len(rs)) + "/" + cmn.IntToString(dataLen)
		rate := fmt.Sprintf("(%.1f%%)", float64(dataLen-len(rs))/float64(dataLen)*100)
		log.Println("[Info] filter splats by alpha:", alpha, ", removed"+rate+":", removed)
	}

	return rs
}
