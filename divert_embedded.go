//go:build windows && (divert_rsrc || divert_embed) && (amd64 || 386 || arm64)

package divert

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/imgk/divert-go/memmod"
)

var (
	winDivert              = (*lazyDLL)(nil)
	winDivertOpen          = (*lazyProc)(nil)
	winDivertCalcChecksums = (*lazyProc)(nil)
)

// Open is ...
func Open(filter string, layer Layer, priority int16, flags uint64) (h *Handle, err error) {
	once.Do(func() {
		dll := newLazyDLL("WinDivert.dll", nil)
		if er := dll.Load(); er != nil {
			err = er
			return
		}
		winDivert = dll

		proc := winDivert.NewProc("WinDivertOpen")
		if er := proc.Find(); er != nil {
			err = er
			return
		}
		winDivertOpen = proc

		proc = winDivert.NewProc("WinDivertHelperCalcChecksums")
		if er := proc.Find(); er != nil {
			err = er
			return
		}
		winDivertCalcChecksums = proc

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

type lazyDLL struct {
	Name   string
	Base   windows.Handle
	mu     sync.Mutex
	module *memmod.Module
	onLoad func(d *lazyDLL)
}

func (p *lazyProc) nameToAddr() (uintptr, error) {
	return p.dll.module.ProcAddressByName(p.Name)
}

func newLazyDLL(name string, onLoad func(d *lazyDLL)) *lazyDLL {
	return &lazyDLL{Name: name, onLoad: onLoad}
}

func (d *lazyDLL) NewProc(name string) *lazyProc {
	return &lazyProc{dll: d, Name: name}
}

type lazyProc struct {
	Name string
	mu   sync.Mutex
	dll  *lazyDLL
	addr uintptr
}

func (p *lazyProc) Find() error {
	if atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&p.addr))) != nil {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.addr != 0 {
		return nil
	}

	err := p.dll.Load()
	if err != nil {
		return fmt.Errorf("Error loading %v DLL: %w", p.dll.Name, err)
	}
	addr, err := p.nameToAddr()
	if err != nil {
		return fmt.Errorf("Error getting %v address: %w", p.Name, err)
	}

	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.addr)), unsafe.Pointer(addr))
	return nil
}

func (p *lazyProc) Addr() uintptr {
	err := p.Find()
	if err != nil {
		panic(err)
	}
	return p.addr
}

func (p *lazyProc) Call(a ...uintptr) (r1, r2 uintptr, lastErr error) {
	return syscall.SyscallN(p.Addr(), a...)
}
