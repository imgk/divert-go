//go:build windows && divert_rsrc && (amd64 || 386 || arm64)

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

	const ourModule windows.Handle = 0
	resInfo, err := windows.FindResource(ourModule, d.Name, windows.RT_RCDATA)
	if err != nil {
		return fmt.Errorf("Unable to find \"%v\" RCDATA resource: %w", d.Name, err)
	}
	data, err := windows.LoadResourceData(ourModule, resInfo)
	if err != nil {
		return fmt.Errorf("Unable to load resource: %w", err)
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
