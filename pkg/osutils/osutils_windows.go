//go:build windows

package osutils

import (
	"os"
	"runtime"
	"syscall"
)

// Убивание процесса в Windows
func TerminateProcess(Pid uint32) error {
	h, e := syscall.OpenProcess(syscall.PROCESS_TERMINATE, false, Pid)
	if e != nil {
		return os.NewSyscallError("OpenProcess", e)
	}
	defer syscall.CloseHandle(h)

	var exitCode uint32 = 0
	e = syscall.TerminateProcess(h, exitCode)
	if e != nil {
		return os.NewSyscallError("TerminateProcess", e)
	}

	runtime.KeepAlive(Pid)
	return nil
}
