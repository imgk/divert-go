// +build windows,!divert_cgo

package divert

type ctlCode uint32

const (
	_METHOD_IN_DIRECT  = 1
	_METHOD_OUT_DIRECT = 2
)

const (
	_FILE_READ_DATA  = 1
	_FILE_WRITE_DATA = 2
)

const (
	_FILE_DEVICE_NETWORK = 0x00000012
)

func _CTL_CODE(DeviceType, Function, Method, Access uint32) ctlCode {
	return ctlCode(((DeviceType) << 16) | ((Access) << 14) | ((Function) << 2) | (Method))
}

var (
	ioCtlInitialize = _CTL_CODE(_FILE_DEVICE_NETWORK, 0x921, _METHOD_OUT_DIRECT, _FILE_READ_DATA|_FILE_WRITE_DATA)
	ioCtlStartup    = _CTL_CODE(_FILE_DEVICE_NETWORK, 0x922, _METHOD_IN_DIRECT, _FILE_READ_DATA|_FILE_WRITE_DATA)
	ioCtlRecv       = _CTL_CODE(_FILE_DEVICE_NETWORK, 0x923, _METHOD_OUT_DIRECT, _FILE_READ_DATA)
	ioCtlSend       = _CTL_CODE(_FILE_DEVICE_NETWORK, 0x924, _METHOD_IN_DIRECT, _FILE_READ_DATA|_FILE_WRITE_DATA)
	ioCtlSetParam   = _CTL_CODE(_FILE_DEVICE_NETWORK, 0x925, _METHOD_IN_DIRECT, _FILE_READ_DATA|_FILE_WRITE_DATA)
	ioCtlGetParam   = _CTL_CODE(_FILE_DEVICE_NETWORK, 0x926, _METHOD_OUT_DIRECT, _FILE_READ_DATA)
	ioCtlShutdown   = _CTL_CODE(_FILE_DEVICE_NETWORK, 0x927, _METHOD_IN_DIRECT, _FILE_READ_DATA|_FILE_WRITE_DATA)
)

func (c ctlCode) String() string {
	switch c {
	case ioCtlInitialize:
		return "IOCTL_WINDIVERT_INITIALIZE"
	case ioCtlStartup:
		return "IOCTL_WINDIVERT_STARTUP"
	case ioCtlRecv:
		return "IOCTL_WINDIVERT_RECV"
	case ioCtlSend:
		return "IOCTL_WINDIVERT_SEND"
	case ioCtlSetParam:
		return "IOCTL_WINDIVERT_SET_PARAM"
	case ioCtlGetParam:
		return "IOCTL_WINDIVERT_GET_PARAM"
	case ioCtlShutdown:
		return "IOCTL_WINDIVERT_SHUTDOWN"
	default:
		return ""
	}
}

type ioCtl struct {
	b1, b2, b3, b4 uint32
}

type recv struct {
	Addr       uint64
	AddrLenPtr uint64
}

type send struct {
	Addr    uint64
	AddrLen uint64
}

type initialize struct {
	Layer    uint32
	Priority uint32
	Flags    uint64
}

type startup struct {
	Flags uint64
	_     uint64
}

type shutdown struct {
	How uint32
	_   uint32
	_   uint64
}

type getParam struct {
	Param uint32
	_     uint32
	Value uint64
}

type setParam struct {
	Value uint64
	Param uint32
	_     uint32
}
