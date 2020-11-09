// +build windows,divert_cgo

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

	"golang.org/x/sys/windows"
)

func Open(filter string, layer Layer, priority int16, flags uint64) (h *Handle, err error) {
	once.Do(func() {
		if er := checkForWow64(); er != nil {
			err = er
			return
		}

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
		return nil, windows.Errno(C.GetLastError())
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
