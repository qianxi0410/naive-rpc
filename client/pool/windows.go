//go:build windows

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
		// BUG: will cause error: The parameter is incorrect.
		n, err := syscall.Read(syscall.Handle(fd), one)

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

// func newOverlapped() (*windows.Overlapped, error) {
// 	var ov windows.Overlapped
// 	h, err := windows.CreateEvent(nil, 0, 1, nil)
// 	if h == 0 {
// 		return nil, err
// 	}
// 	ov.HEvent = h
// 	return &ov, nil
// }

// func readClosed(conn net.Conn) bool {
// 	var readClosed bool
// 	f := func(fd uintptr) bool {
// 		one := []byte{0}
// 		var n uint32
// 		ov, err := newOverlapped()
// 		if err != nil {
// 			return false
// 		}
// 		defer windows.CloseHandle(ov.HEvent)

// 		err = windows.ReadFile(windows.Handle(fd), one, &n, ov)
// 		if err != nil && err == syscall.ERROR_IO_PENDING {
// 			if err = windows.GetOverlappedResult(windows.Handle(fd), ov, &n, true); err != nil {
// 				return false
// 			}
// 		}

// 		// if err != nil && err != syscall.EAGAIN || err == nil && n == 0 {
// 		// 	readClosed = true
// 		// }

// 		return true
// 	}

// 	conn.SetReadDeadline(time.Time{})

// 	var rawConn syscall.RawConn
// 	var err error
// 	switch v := conn.(type) {
// 	case *net.TCPConn:
// 		rawConn, err = v.SyscallConn()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	default:
// 		return false
// 	}

// 	// half closed
// 	if err := rawConn.Read(f); err != nil {
// 		return true
// 	}

// 	return readClosed
// }
