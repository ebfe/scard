// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs -- -I/usr/include/PCSC const.go

package scard

type Attrib uint32

const (
	AttrVendorName           Attrib = 0x10100
	AttrVendorIfdType        Attrib = 0x10101
	AttrVendorIfdVersion     Attrib = 0x10102
	AttrVendorIfdSerialNo    Attrib = 0x10103
	AttrChannelId            Attrib = 0x20110
	AttrAsyncProtocolTypes   Attrib = 0x30120
	AttrDefaultClk           Attrib = 0x30121
	AttrMaxClk               Attrib = 0x30122
	AttrDefaultDataRate      Attrib = 0x30123
	AttrMaxDataRate          Attrib = 0x30124
	AttrMaxIfsd              Attrib = 0x30125
	AttrSyncProtocolTypes    Attrib = 0x30126
	AttrPowerMgmtSupport     Attrib = 0x40131
	AttrUserToCardAuthDevice Attrib = 0x50140
	AttrUserAuthInputDevice  Attrib = 0x50142
	AttrCharacteristics      Attrib = 0x60150
	AttrCurrentProtocolType  Attrib = 0x80201
	AttrCurrentClk           Attrib = 0x80202
	AttrCurrentF             Attrib = 0x80203
	AttrCurrentD             Attrib = 0x80204
	AttrCurrentN             Attrib = 0x80205
	AttrCurrentW             Attrib = 0x80206
	AttrCurrentIfsc          Attrib = 0x80207
	AttrCurrentIfsd          Attrib = 0x80208
	AttrCurrentBwt           Attrib = 0x80209
	AttrCurrentCwt           Attrib = 0x8020a
	AttrCurrentEbcEncoding   Attrib = 0x8020b
	AttrExtendedBwt          Attrib = 0x8020c
	AttrIccPresence          Attrib = 0x90300
	AttrIccInterfaceStatus   Attrib = 0x90301
	AttrCurrentIoState       Attrib = 0x90302
	AttrAtrString            Attrib = 0x90303
	AttrIccTypePerAtr        Attrib = 0x90304
	AttrEscReset             Attrib = 0x7a000
	AttrEscCancel            Attrib = 0x7a003
	AttrEscAuthrequest       Attrib = 0x7a005
	AttrMaxinput             Attrib = 0x7a007
	AttrDeviceUnit           Attrib = 0x7fff0001
	AttrDeviceInUse          Attrib = 0x7fff0002
	AttrDeviceFriendlyName   Attrib = 0x7fff0003
	AttrDeviceSystemName     Attrib = 0x7fff0004
	AttrSupressT1IfsRequest  Attrib = 0x7fff0007
)

const (
	ErrSuccess                scardError = 0x0
	ErrInternalError          scardError = 0x80100001
	ErrCancelled              scardError = 0x80100002
	ErrInvalidHandle          scardError = 0x80100003
	ErrInvalidParameter       scardError = 0x80100004
	ErrInvalidTarget          scardError = 0x80100005
	ErrNoMemory               scardError = 0x80100006
	ErrWaitedTooLong          scardError = 0x80100007
	ErrInsufficientBuffer     scardError = 0x80100008
	ErrUnknownReader          scardError = 0x80100009
	ErrTimeout                scardError = 0x8010000a
	ErrSharingViolation       scardError = 0x8010000b
	ErrNoSmartcard            scardError = 0x8010000c
	ErrUnknownCard            scardError = 0x8010000d
	ErrCantDispose            scardError = 0x8010000e
	ErrProtoMismatch          scardError = 0x8010000f
	ErrNotReady               scardError = 0x80100010
	ErrInvalidValue           scardError = 0x80100011
	ErrSystemCancelled        scardError = 0x80100012
	ErrCommError              scardError = 0x80100013
	ErrUnknownError           scardError = 0x80100014
	ErrInvalidAtr             scardError = 0x80100015
	ErrNotTransacted          scardError = 0x80100016
	ErrReaderUnavailable      scardError = 0x80100017
	ErrShutdown               scardError = 0x80100018
	ErrPciTooSmall            scardError = 0x80100019
	ErrReaderUnsupported      scardError = 0x8010001a
	ErrDuplicateReader        scardError = 0x8010001b
	ErrCardUnsupported        scardError = 0x8010001c
	ErrNoService              scardError = 0x8010001d
	ErrServiceStopped         scardError = 0x8010001e
	ErrUnexpected             scardError = 0x8010001f
	ErrUnsupportedFeature     scardError = 0x8010001f
	ErrIccInstallation        scardError = 0x80100020
	ErrIccCreateorder         scardError = 0x80100021
	ErrFileNotFound           scardError = 0x80100024
	ErrNoDir                  scardError = 0x80100025
	ErrNoFile                 scardError = 0x80100026
	ErrNoAccess               scardError = 0x80100027
	ErrWriteTooMany           scardError = 0x80100028
	ErrBadSeek                scardError = 0x80100029
	ErrInvalidChv             scardError = 0x8010002a
	ErrUnknownResMng          scardError = 0x8010002b
	ErrNoSuchCertificate      scardError = 0x8010002c
	ErrCertificateUnavailable scardError = 0x8010002d
	ErrNoReadersAvailable     scardError = 0x8010002e
	ErrCommDataLost           scardError = 0x8010002f
	ErrNoKeyContainer         scardError = 0x80100030
	ErrServerTooBusy          scardError = 0x80100031
	ErrUnsupportedCard        scardError = 0x80100065
	ErrUnresponsiveCard       scardError = 0x80100066
	ErrUnpoweredCard          scardError = 0x80100067
	ErrResetCard              scardError = 0x80100068
	ErrRemovedCard            scardError = 0x80100069
	ErrSecurityViolation      scardError = 0x8010006a
	ErrWrongChv               scardError = 0x8010006b
	ErrChvBlocked             scardError = 0x8010006c
	ErrEof                    scardError = 0x8010006d
	ErrCancelledByUser        scardError = 0x8010006e
	ErrCardNotAuthenticated   scardError = 0x8010006f
)

type Protocol uint32

const (
	ProtocolUndefined Protocol = 0x0
	ProtocolT0        Protocol = 0x1
	ProtocolT1        Protocol = 0x2
	ProtocolAny       Protocol = ProtocolT0 | ProtocolT1
)

type ShareMode uint32

const (
	ShareExclusive ShareMode = 0x1
	ShareShared    ShareMode = 0x2
	ShareDirect    ShareMode = 0x3
)

type Disposition uint32

const (
	LeaveCard   Disposition = 0x0
	ResetCard   Disposition = 0x1
	UnpowerCard Disposition = 0x2
	EjectCard   Disposition = 0x3
)

type State uint32

const (
	Unknown    State = 0x1
	Absent     State = 0x2
	Present    State = 0x4
	Swallowed  State = 0x8
	Powered    State = 0x10
	Negotiable State = 0x20
	Specific   State = 0x40
)

type StateFlag uint32

const (
	StateUnaware     StateFlag = 0x0
	StateIgnore      StateFlag = 0x1
	StateChanged     StateFlag = 0x2
	StateUnknown     StateFlag = 0x4
	StateUnavailable StateFlag = 0x8
	StateEmpty       StateFlag = 0x10
	StatePresent     StateFlag = 0x20
	StateAtrmatch    StateFlag = 0x40
	StateExclusive   StateFlag = 0x80
	StateInuse       StateFlag = 0x100
	StateMute        StateFlag = 0x200
	StateUnpowered   StateFlag = 0x400
)

const (
	maxBufferSize         = 0x108
	maxBufferSizeExtended = 0x1000c
	maxReadername         = 0x80
	maxAtrSize            = 0x21
)

const (
	infiniteTimeout = 0xffffffff
)
