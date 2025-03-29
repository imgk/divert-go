//go:build windows && !divert_cgo && !divert_rsrc && !divert_embed && (amd64 || 386 || arm64)

package divert

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/sys/windows"
)

var (
	winDivert              = (*windows.DLL)(nil)
	winDivertOpen          = (*windows.Proc)(nil)
	winDivertCalcChecksums = (*windows.Proc)(nil)
)

// Open is ...
func Open(filter string, layer Layer, priority int16, flags uint64) (h *Handle, err error) {
	once.Do(func() {
		dll, er := windows.LoadDLL("WinDivert.dll")
		if er != nil {
			err = er
			return
		}
		winDivert = dll

		proc, er := winDivert.FindProc("WinDivertOpen")
		if er != nil {
			err = er
			return
		}
		winDivertOpen = proc

		proc, er = winDivert.FindProc("WinDivertHelperCalcChecksums")
		if er != nil {
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
