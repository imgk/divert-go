//go:build windows && divert_cgo && (amd64 || 386 || arm64)

package divert

// #cgo CFLAGS: -I${SRCDIR}/divert -Wno-incompatible-pointer-types
// #define WINDIVERTEXPORT static
// #include "windivert.c"
import "C"

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Open is ...
func Open(filter string, layer Layer, priority int16, flags uint64) (h *Handle, err error) {
	once.Do(func() {
		vers := map[string]struct{}{
			"2.0": {},
			"2.1": {},
			"2.2": {},
		}
		ver, er := func() (ver string, err error) {
			h, err := open("false", LayerNetwork, PriorityDefault, FlagDefault)
			if err != nil {
				return
			}
			defer func() {
				err = h.Close()
			}()

			major, err := h.GetParam(VersionMajor)
			if err != nil {
				return
			}

			minor, err := h.GetParam(VersionMinor)
			if err != nil {
				return
			}

			ver = strings.Join([]string{strconv.Itoa(int(major)), strconv.Itoa(int(minor))}, ".")
			return
		}()
		if er != nil {
			err = er
			return
		}
		if _, ok := vers[ver]; !ok {
			err = fmt.Errorf("unsupported windivert version: %v", ver)
		}
	})
	if err != nil {
		return
	}

	return open(filter, layer, priority, flags)
}

func open(filter string, layer Layer, priority int16, flags uint64) (h *Handle, err error) {
	if priority < PriorityLowest || priority > PriorityHighest {
		return nil, errPriority
	}

	runtime.LockOSThread()
	hd := C.WinDivertOpen(C.CString(filter), C.WINDIVERT_LAYER(layer), C.int16_t(priority), C.uint64_t(flags))
	runtime.UnlockOSThread()

	if hd == C.HANDLE(C.INVALID_HANDLE_VALUE) {
		return nil, Error(C.GetLastError())
	}

	rEvent, _ := windows.CreateEvent(nil, 0, 0, nil)
	wEvent, _ := windows.CreateEvent(nil, 0, 0, nil)

	return &Handle{
		Mutex:  sync.Mutex{},
		Handle: windows.Handle(uintptr(unsafe.Pointer(hd))),
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
	re := C.WinDivertHelperCalcChecksums(unsafe.Pointer(&buffer[0]), C.UINT(len(buffer)), (*C.WINDIVERT_ADDRESS)(unsafe.Pointer(address)), C.uint64_t(flags))
	return re == C.TRUE
}
