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
	"unsafe"

	"golang.org/x/sys/windows"
)

var once = sync.Once{}

func checkVersion() error {
	if err := checkForWow64(); err != nil {
		return err
	}

	vers := map[string]struct{}{
		"2.0": struct{}{},
		"2.1": struct{}{},
		"2.2": struct{}{},
	}
	ver, err := GetVersion()
	if err != nil {
		return err
	}
	if _, ok := vers[ver]; !ok {
		return fmt.Errorf("unsupported windivert version: %v", ver)
	}
	return nil
}

// Get version info of windivert
func GetVersion() (ver string, err error) {
	h, err := Open("false", LayerNetwork, PriorityDefault, FlagDefault)
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
}

func checkForWow64() error {
	var b bool
	err := windows.IsWow64Process(windows.CurrentProcess(), &b)
	if err != nil {
		return fmt.Errorf("Unable to determine whether the process is running under WOW64: %v", err)
	}
	if b {
		return fmt.Errorf("You must use the 64-bit version of WireGuard on this computer.")
	}
	return nil
}

func getLastError() error {
	if errno := windows.Errno(C.GetLastError()); errno != windows.ERROR_SUCCESS {
		return Error(errno)
	}
	return nil
}

type Handle struct {
	sync.Mutex
	windows.Handle
}

func Open(filter string, layer Layer, priority int16, flags uint64) (h *Handle, err error) {
	once.Do(func() {
		err = checkVersion()
	})
	if err != nil {
		return
	}

	if priority < PriorityLowest || priority > PriorityHighest {
		return nil, errPriority
	}

	runtime.LockOSThread()
	hd := C.WinDivertOpen(C.CString(filter), C.WINDIVERT_LAYER(layer), C.int16_t(priority), C.uint64_t(flags))
	runtime.UnlockOSThread()

	if hd == C.HANDLE(C.INVALID_HANDLE_VALUE) {
		return nil, getLastError()
	}

	return &Handle{
		Mutex:  sync.Mutex{},
		Handle: windows.Handle(hd),
	}, nil
}

func (h *Handle) Recv(buffer []byte, address *Address) (uint, error) {
	recvLen := uint(0)

	b := C.WinDivertRecv(C.HANDLE(h.Handle), unsafe.Pointer(&buffer[0]), C.uint(len(buffer)), (*C.uint)(unsafe.Pointer(&recvLen)), C.PWINDIVERT_ADDRESS(unsafe.Pointer(address)))
	if b == C.FALSE {
		return 0, getLastError()
	}

	return recvLen, nil
}

func (h *Handle) RecvEx(buffer []byte, address []Address, overlapped *windows.Overlapped) (uint, uint, error) {
	recvLen := uint(0)

	addrLen := uint(len(address)) * uint(unsafe.Sizeof(C.WINDIVERT_ADDRESS{}))
	b := C.WinDivertRecvEx(C.HANDLE(h.Handle), unsafe.Pointer(&buffer[0]), C.uint(len(buffer)), (*C.uint)(unsafe.Pointer(&recvLen)), C.uint64_t(0), C.PWINDIVERT_ADDRESS(unsafe.Pointer(&address[0])), (*C.uint)(unsafe.Pointer(&addrLen)), C.LPOVERLAPPED(unsafe.Pointer(overlapped)))
	if b == C.FALSE {
		return 0, 0, getLastError()
	}
	addrLen /= uint(unsafe.Sizeof(C.WINDIVERT_ADDRESS{}))

	return recvLen, addrLen, nil
}

func (h *Handle) Send(buffer []byte, address *Address) (uint, error) {
	sendLen := uint(0)

	b := C.WinDivertSend(C.HANDLE(h.Handle), unsafe.Pointer(&buffer[0]), C.uint(len(buffer)), (*C.uint)(unsafe.Pointer(&sendLen)), (*C.WINDIVERT_ADDRESS)(unsafe.Pointer(address)))
	if b == C.FALSE {
		return 0, getLastError()
	}

	return sendLen, nil
}

func (h *Handle) SendEx(buffer []byte, address []Address, overlapped *windows.Overlapped) (uint, error) {
	sendLen := uint(0)

	b := C.WinDivertSendEx(C.HANDLE(h.Handle), unsafe.Pointer(&buffer[0]), C.uint(len(buffer)), (*C.uint)(unsafe.Pointer(&sendLen)), C.uint64_t(0), (*C.WINDIVERT_ADDRESS)(unsafe.Pointer(&address[0])), C.uint(uint(len(address))*uint(unsafe.Sizeof(C.WINDIVERT_ADDRESS{}))), C.LPOVERLAPPED(unsafe.Pointer(overlapped)))
	if b == C.FALSE {
		return 0, getLastError()
	}

	return sendLen, nil
}

func (h *Handle) Shutdown(how Shutdown) error {
	b := C.WinDivertShutdown(C.HANDLE(h.Handle), C.WINDIVERT_SHUTDOWN(how))
	if b == C.FALSE {
		return getLastError()
	}

	return nil
}

func (h *Handle) Close() error {
	b := C.WinDivertClose(C.HANDLE(h.Handle))
	if b == C.FALSE {
		return getLastError()
	}

	return nil
}

func (h *Handle) GetParam(p Param) (uint64, error) {
	v := uint64(0)

	b := C.WinDivertGetParam(C.HANDLE(h.Handle), C.WINDIVERT_PARAM(p), (*C.uint64_t)(unsafe.Pointer(&v)))
	if b == C.FALSE {
		err := getLastError()
		return v, err
	}

	return v, nil
}

func (h *Handle) SetParam(p Param, v uint64) error {
	switch p {
	case QueueLength:
		if v < QueueLengthMin || v > QueueLengthMax {
			return errQueueLength
		}
	case QueueTime:
		if v < QueueTimeMin || v > QueueTimeMax {
			return errQueueTime
		}
	case QueueSize:
		if v < QueueSizeMin || v > QueueSizeMax {
			return errQueueSize
		}
	default:
		return errQueueParam
	}

	b := C.WinDivertSetParam(C.HANDLE(h.Handle), C.WINDIVERT_PARAM(p), C.uint64_t(v))
	if b == C.FALSE {
		return getLastError()
	}

	return nil
}
