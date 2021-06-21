// +build windows

package divert

import "unsafe"

// Ethernet is ...
type Ethernet struct {
	InterfaceIndex    uint32
	SubInterfaceIndex uint32
	_                 [7]uint64
}

// Network is ...
// The WINDIVERT_LAYER_NETWORK and WINDIVERT_LAYER_NETWORK_FORWARD layers allow the user
// application to capture/block/inject network packets passing to/from (and through) the
// local machine. Due to technical limitations, process ID information is not available
// at these layers.
type Network struct {
	InterfaceIndex    uint32
	SubInterfaceIndex uint32
	_                 [7]uint64
}

// Socket is ...
// The WINDIVERT_LAYER_SOCKET layer can capture or block events corresponding to socket
// operations, such as bind(), connect(), listen(), etc., or the termination of socket
// operations, such as a TCP socket disconnection. Unlike the flow layer, most socket-related
// events can be blocked. However, it is not possible to inject new or modified socket events.
// Process ID information (of the process responsible for the socket operation) is available
// at this layer. Due to technical limitations, this layer cannot capture events that occurred
// before the handle was opened.
type Socket struct {
	EndpointID       uint64
	ParentEndpointID uint64
	ProcessID        uint32
	LocalAddress     [16]uint8
	RemoteAddress    [16]uint8
	LocalPort        uint16
	RemotePort       uint16
	Protocol         uint8
	_                [3]uint8
	_                uint32
}

// Flow is ...
// The WINDIVERT_LAYER_FLOW layer captures information about network flow establishment/deletion
// events. Here, a flow represents either (1) a TCP connection, or (2) an implicit "flow" created
// by the first sent/received packet for non-TCP traffic, e.g., UDP. Old flows are deleted when
// the corresponding connection is closed (for TCP), or based on an activity timeout (non-TCP).
// Flow-related events can be captured, but not blocked nor injected. Process ID information is
// also available at this layer. Due to technical limitations, the WINDIVERT_LAYER_FLOW layer
// cannot capture flow events that occurred before the handle was opened.
type Flow struct {
	EndpointID       uint64
	ParentEndpointID uint64
	ProcessID        uint32
	LocalAddress     [16]uint8
	RemoteAddress    [16]uint8
	LocalPort        uint16
	RemotePort       uint16
	Protocol         uint8
	_                [3]uint8
	_                uint32
}

// Reflect is ...
// Finally, the WINDIVERT_LAYER_REFLECT layer can capture events relating to WinDivert itself,
// such as when another process opens a new WinDivert handle, or closes an old WinDivert handle.
// WinDivert events can be captured but not injected nor blocked. Process ID information
// (of the process responsible for opening the WinDivert handle) is available at this layer.
// This layer also returns data in the form of an "object" representation of the filter string
// used to open the handle. The object representation can be converted back into a human-readable
// filter string using the WinDivertHelperFormatFilter() function. This layer can also capture
// events that occurred before the handle was opened. This layer cannot capture events related
// to other WINDIVERT_LAYER_REFLECT-layer handles.
type Reflect struct {
	TimeStamp int64
	ProcessID uint32
	layer     uint32
	Flags     uint64
	Priority  int16
	_         int16
	_         int32
	_         [4]uint64
}

// Layer is ...
func (r *Reflect) Layer() Layer {
	return Layer(r.layer)
}

// Address is ...
type Address struct {
	Timestamp int64
	layer     uint8
	event     uint8
	Flags     uint8
	_         uint8
	length    uint32
	union     [64]uint8
}

// Layer is ...
func (a *Address) Layer() Layer {
	return Layer(a.layer)
}

// SetLayer is ...
func (a *Address) SetLayer(layer Layer) {
	a.layer = uint8(layer)
}

// Event is ...
func (a *Address) Event() Event {
	return Event(a.event)
}

// SetEvent is ...
func (a *Address) SetEvent(event Event) {
	a.event = uint8(event)
}

// Length is ...
func (a *Address) Length() uint32 {
	return a.length >> 12
}

// SetLength is ...
func (a *Address) SetLength(n uint32) {
	a.length = n << 12
}

// Ethernet is ...
func (a *Address) Ethernet() *Ethernet {
	return (*Ethernet)(unsafe.Pointer(&a.union))
}

// Network is ...
func (a *Address) Network() *Network {
	return (*Network)(unsafe.Pointer(&a.union))
}

// Socket is ...
func (a *Address) Socket() *Socket {
	return (*Socket)(unsafe.Pointer(&a.union))
}

// Flow is ...
func (a *Address) Flow() *Flow {
	return (*Flow)(unsafe.Pointer(&a.union))
}

// Reflect is ...
func (a *Address) Reflect() *Reflect {
	return (*Reflect)(unsafe.Pointer(&a.union))
}
