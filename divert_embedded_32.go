//go:build windows && divert_embedded && (386 || arm)
// +build windows
// +build divert_embedded
// +build 386 arm

package divert

import (
	"runtime"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

func open(filter string, layer Layer, priority int16, flags uint64) (h *Handle, err error) {
	if priority < PriorityLowest || priority > PriorityHighest {
		return nil, errPriority
	}

	filterPtr, err := windows.BytePtrFromString(filter)
	if err != nil {
		return nil, err
	}

	runtime.LockOSThread()
	// 386: Shadow panics on Windows 32-bit system, as `flags` is `uint64` while `uintptr` is `uint32` for 32-bit system.
	hd, _, err := winDivertOpen.Call(uintptr(unsafe.Pointer(filterPtr)), uintptr(layer), uintptr(priority), uintptr(flags), 0)
	runtime.UnlockOSThread()

	if windows.Handle(hd) == windows.InvalidHandle {
		return nil, Error(err.(windows.Errno))
	}

	rEvent, _ := windows.CreateEvent(nil, 0, 0, nil)
	wEvent, _ := windows.CreateEvent(nil, 0, 0, nil)

	return &Handle{
		Mutex:  sync.Mutex{},
		Handle: windows.Handle(hd),
		rOverlapped: windows.Overlapped{
			HEvent: rEvent,
		},
		wOverlapped: windows.Overlapped{
			HEvent: wEvent,
		},
	}, nil
}

// CalcChecksums is ...
func CalcChecksums(buffer []byte, address *Address, flags uint64) bool {
	re, _, err := winDivertCalcChecksums.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)), uintptr(unsafe.Pointer(address)), uintptr(flags), 0)
	if err != nil {
		return false
	}
	return re != 0
}
