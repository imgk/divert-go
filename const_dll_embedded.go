// +build windows,!divert_cgo windows,divert_embedded

package divert

const (
	LayerNetwork        Layer = 0
	LayerNetworkForward Layer = 1
	LayerFlow           Layer = 2
	LayerSocket         Layer = 3
	LayerReflect        Layer = 4
	//LayerEthernet       Layer = 5
)

const (
	EventNetworkPacket   Event = 0
	EventFlowEstablished Event = 1
	EventFlowDeleted     Event = 2
	EventSocketBind      Event = 3
	EventSocketConnect   Event = 4
	EventSocketListen    Event = 5
	EventSocketAccept    Event = 6
	EventSocketClose     Event = 7
	EventReflectOpen     Event = 8
	EventReflectClose    Event = 9
	//EventEthernetFrame   Event = 10
)

const (
	ShutdownRecv Shutdown = 0
	ShutdownSend Shutdown = 1
	ShutdownBoth Shutdown = 2
)

const (
	QueueLength  Param = 0
	QueueTime    Param = 1
	QueueSize    Param = 2
	VersionMajor Param = 3
	VersionMinor Param = 4
)

const (
	FlagDefault   = 0x0000
	FlagSniff     = 0x0001
	FlagDrop      = 0x0002
	FlagRecvOnly  = 0x0004
	FlagSendOnly  = 0x0008
	FlagNoInstall = 0x0010
	FlagFragments = 0x0020
)

const (
	PriorityDefault    = 0
	PriorityHighest    = 3000
	PriorityLowest     = -3000
	QueueLengthDefault = 4096
	QueueLengthMin     = 32
	QueueLengthMax     = 16384
	QueueTimeDefault   = 2000
	QueueTimeMin       = 100
	QueueTimeMax       = 16000
	QueueSizeDefault   = 4194304
	QueueSizeMin       = 65535
	QueueSizeMax       = 33554432
)

const (
	ChecksumDefault  = 0
	NoIPChecksum     = 1
	NoICMPChekcsum   = 2
	NoICMPV6Checksum = 4
	NoTCPChekcsum    = 8
	NoUDPChecksum    = 16
)

const (
	BatchMax = 0xff
	MTUMax   = 40 + 0xffff
)
