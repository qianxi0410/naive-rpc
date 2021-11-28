//go:build !windows

package pool

import (
	"net"
	"syscall"
	"time"
)

func readClosed(conn net.Conn) bool {
	var readClosed bool
	f := func(fd uintptr) bool {
		one := []byte{0}
		n, err := syscall.Read(int(fd), one)

		if err != nil && err != syscall.EAGAIN || err == nil && n == 0 {
			readClosed = true
		}

		return true
	}

	conn.SetReadDeadline(time.Time{})

	var rawConn syscall.RawConn
	switch v := conn.(type) {
	case *net.TCPConn:
		rawConn, _ = v.SyscallConn()
	case *net.UnixConn:
		rawConn, _ = v.SyscallConn()
	default:
		return false
	}

	// half closed
	if err := rawConn.Read(f); err != nil {
		return true
	}

	return readClosed
}
