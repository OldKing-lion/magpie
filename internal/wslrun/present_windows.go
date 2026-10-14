//go:build windows

package wslrun

import "golang.org/x/sys/windows/registry"

// lxssHasDistro is whether Windows has registered a WSL distro. The Lxss
// key exists on machines that have never installed one; only a default or
// a DistributionName means wsl.exe will list something instead of asking
// to install.
func lxssHasDistro() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Lxss`, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return false
	}
	defer k.Close()
	if v, _, err := k.GetStringValue("DefaultDistribution"); err == nil && v != "" {
		return true
	}
	names, _ := k.ReadSubKeyNames(-1)
	for _, name := range names {
		sub, err := registry.OpenKey(k, name, registry.QUERY_VALUE)
		if err != nil {
			continue
		}
		n, _, err := sub.GetStringValue("DistributionName")
		sub.Close()
		if err == nil && n != "" {
			return true
		}
	}
	return false
}
