//go:build darwin || linux

package runtime

import (
	"errors"
	"golang.org/x/sys/unix"
	"os"
)

func lockChatSettings(f *os.File) error { return unix.Flock(int(f.Fd()), unix.LOCK_EX|unix.LOCK_NB) }
func unlockChatSettings(f *os.File)     { _ = unix.Flock(int(f.Fd()), unix.LOCK_UN) }
func chatLockBusy(err error) bool {
	return errors.Is(err, unix.EWOULDBLOCK) || errors.Is(err, unix.EAGAIN)
}
