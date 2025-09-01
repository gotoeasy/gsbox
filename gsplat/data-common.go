package gsplat

import (
	"gsbox/cmn"
)

func IsOutputSpx() bool {
	if Args.HasCmd("p2x", "ply2spx", "s2x", "splat2spx", "z2x", "spz2spx", "x2x", "spx2spx", "k2x", "ksplat2spx") {
		return true
	}
	if Args.HasCmd("join") && cmn.Endwiths(Args.GetArgIgnorecase("-o", "--output"), ".spx", true) {
		return true
	}
	return false
}

func IsOutputSpz() bool {
	if Args.HasCmd("p2z", "ply2spz", "s2z", "splat2spz", "z2z", "spz2spz", "x2z", "spx2spz", "k2z", "ksplat2spz") {
		return true
	}
	if Args.HasCmd("join") && cmn.Endwiths(Args.GetArgIgnorecase("-o", "--output"), ".spz", true) {
		return true
	}
	return false
}

func IsOutputPly() bool {
	if Args.HasCmd("p2p", "ply2ply", "s2p", "splat2ply", "z2p", "spz2ply", "x2p", "spx2ply", "k2p", "ksplat2ply") {
		return true
	}
	if Args.HasCmd("join") && cmn.Endwiths(Args.GetArgIgnorecase("-o", "--output"), ".ply", true) {
		return true
	}
	return false
}

func IsOutputSplat() bool {
	if Args.HasCmd("p2s", "ply2splat", "s2s", "splat2splat", "z2s", "spz2splat", "x2s", "spx2splat", "k2s", "ksplat2splat") {
		return true
	}
	if Args.HasCmd("join") && cmn.Endwiths(Args.GetArgIgnorecase("-o", "--output"), ".splat", true) {
		return true
	}
	return false
}
