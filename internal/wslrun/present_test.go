package wslrun

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/yetone/magpie/internal/proc"
)

func TestPresentEnvOff(t *testing.T) {
	resetDenied(t)
	for _, v := range []string{"0", "false", "no", "off", "OFF"} {
		t.Setenv("MAGPIE_WSL", v)
		if present() {
			t.Fatalf("MAGPIE_WSL=%s: Present", v)
		}
	}
}

func TestRunDoesNotStartWSLWhenAbsent(t *testing.T) {
	resetDenied(t)
	t.Setenv("MAGPIE_WSL", "0")
	old := start
	t.Cleanup(func() { start = old })
	start = func(ctx context.Context, args ...string) *exec.Cmd {
		t.Errorf("wsl.exe %v", args)
		return proc.CommandContext(ctx, "true")
	}
	out, err := Run(time.Second, "-l", "-q")
	if err == nil {
		t.Fatal("Run succeeded with MAGPIE_WSL=0")
	}
	if out != nil {
		t.Fatalf("out = %q", out)
	}
}

func TestInstallStub(t *testing.T) {
	zh := "适用于 Linux 的 Windows 子系统必须更新到最新版本才能继续。可通过运行 “wsl.exe --update” 进行更新。\n" +
		"有关详细信息，请访问 https://aka.ms/wslinstall\n" +
		"按任意键安装适用于 Linux 的 Windows 子系统。按 ESC 或 CTRL-C 取消。此提示将在 60 秒后超时。\n"
	en := "Windows Subsystem for Linux must be updated to the latest version to continue. Please run \"wsl.exe --update\" to update.\n" +
		"For more information, visit https://aka.ms/wslinstall\n" +
		"Press any key to install Windows Subsystem for Linux. Press ESC or CTRL-C to cancel. This prompt will time out in 60 seconds.\n"
	if !installStub([]byte(zh)) {
		t.Fatal("chinese install prompt")
	}
	if !installStub([]byte(en)) {
		t.Fatal("english install prompt")
	}
	if !installStub(utf16le(zh)) {
		t.Fatal("utf-16 chinese install prompt")
	}
	if installStub([]byte("Ubuntu\r\nDebian\r\n")) {
		t.Fatal("a distro list taken for the install prompt")
	}
	if installStub(nil) {
		t.Fatal("empty output taken for the install prompt")
	}
}

func TestInstallStubStopsLaterRuns(t *testing.T) {
	resetDenied(t)
	deny()
	if present() {
		t.Fatal("still Present after the inbox prompt")
	}
	old := start
	t.Cleanup(func() { start = old })
	start = func(ctx context.Context, args ...string) *exec.Cmd {
		t.Errorf("wsl.exe %v", args)
		return proc.CommandContext(ctx, "true")
	}
	if _, err := Run(time.Second, "-l", "-q"); err == nil {
		t.Fatal("Run after the inbox prompt")
	}
}

func resetDenied(t *testing.T) {
	t.Helper()
	denyMu.Lock()
	denied = false
	denyMu.Unlock()
	t.Cleanup(func() {
		denyMu.Lock()
		denied = false
		denyMu.Unlock()
	})
}
