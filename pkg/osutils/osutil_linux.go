//go:build linux

package osutils

import (
	"golang.org/x/sys/unix"
	"syscall"
)

// Убивание процесса в linux
func TerminateProcess(Pid uint32) error {
	var signal syscall.Signal = 1
	err := unix.Kill(int(Pid), signal)
	return err
}
