package wslrun

import (
	"os"
	"runtime"
	"strings"
	"sync"
)

// Present is whether magpie should talk to wsl.exe. The inbox wsl.exe on a
// Windows machine with no distro opens a console asking to install WSL
// (and another every few seconds while magpie lists agents), so nothing
// calls it unless a distro is registered. MAGPIE_WSL=0 turns looking off;
// MAGPIE_WSL=1 forces it on. Off Windows it is always false.
func Present() bool { return present() }

var (
	denyMu sync.Mutex
	denied bool
)

func present() bool {
	denyMu.Lock()
	d := denied
	denyMu.Unlock()
	if d {
		return false
	}
	if runtime.GOOS != "windows" {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(os.Getenv("MAGPIE_WSL"))) {
	case "0", "false", "no", "off":
		return false
	case "1", "true", "yes", "on":
		return true
	}
	return lxssHasDistro()
}

// deny remembers that wsl.exe is the inbox installer, so later calls do
// not open another console.
func deny() {
	denyMu.Lock()
	denied = true
	denyMu.Unlock()
}

func installStub(b []byte) bool {
	s := strings.ToLower(decode(b))
	for _, n := range []string{
		"aka.ms/wslinstall",
		"wsl.exe --update",
		"press any key to install",
		"press any key to install windows subsystem",
		"按任意键安装",
		"必须更新到最新版本",
		"windows subsystem for linux must be updated",
	} {
		if strings.Contains(s, strings.ToLower(n)) {
			return true
		}
	}
	return false
}
