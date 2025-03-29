//go:build windows && divert_embed && amd64

package divert

import (
	"fmt"
	"sync/atomic"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/imgk/divert-go/memmod"
)

func (d *lazyDLL) Load() error {
	if atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&d.module))) != nil {
		return nil
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.module != nil {
		return nil
	}

	data, err := f.ReadFile("WinDivert/x64/WinDivert.dll")
	if err != nil {
		return fmt.Errorf("load dll error: %w", err)
	}
	module, err := memmod.LoadLibrary(data)
	if err != nil {
		return fmt.Errorf("Unable to load library: %w", err)
	}
	d.Base = windows.Handle(module.BaseAddr())

	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&d.module)), unsafe.Pointer(module))
	if d.onLoad != nil {
		d.onLoad(d)
	}
	return nil
}
