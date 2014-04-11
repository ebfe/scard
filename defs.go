package scard

type Protocol uint32

const (
	PROTOCOL_UNDEFINED Protocol = 0
	PROTOCOL_T0        Protocol = 1
	PROTOCOL_T1        Protocol = 2
	PROTOCOL_ANY       Protocol = PROTOCOL_T0 | PROTOCOL_T1
)

type ShareMode uint32

const (
	SHARE_EXCLUSIVE ShareMode = 1
	SHARE_SHARED    ShareMode = 2
	SHARE_DIRECT    ShareMode = 3
)

type Disposition uint32

const (
	LEAVE_CARD   Disposition = 0
	RESET_CARD   Disposition = 1
	UNPOWER_CARD Disposition = 2
	EJECT_CARD   Disposition = 3
)

type State uint32

const (
	UNKNOWN    State = 0x0001
	ABSENT     State = 0x0002
	PRESENT    State = 0x0004
	SWALLOWED  State = 0x0008
	POWERED    State = 0x0010
	NEGOTIABLE State = 0x0020
	SPECIFIC   State = 0x0040
)

type StateFlag uint32

const (
	STATE_UNAWARE     StateFlag = 0x0000
	STATE_IGNORE      StateFlag = 0x0001
	STATE_CHANGED     StateFlag = 0x0002
	STATE_UNKNOWN     StateFlag = 0x0004
	STATE_UNAVAILABLE StateFlag = 0x0008
	STATE_EMPTY       StateFlag = 0x0010
	STATE_PRESENT     StateFlag = 0x0020
	STATE_ATRMATCH    StateFlag = 0x0040
	STATE_EXCLUSIVE   StateFlag = 0x0080
	STATE_INUSE       StateFlag = 0x0100
	STATE_MUTE        StateFlag = 0x0200
	STATE_UNPOWERED   StateFlag = 0x0400
)

type Timeout uint32

const (
	INFINITE Timeout = 0xffffffff
)

const (
	MAX_BUFFER_SIZE          = 264
	MAX_BUFFER_SIZE_EXTENDED = 4 + 3 + (1 << 16) + 3 + 2
	MAX_READERNAME           = 128
	MAX_ATR_SIZE             = 33
)

type CardStatus struct {
	Reader         string
	State          State
	ActiveProtocol Protocol
	ATR            []byte
}

type ReaderState struct {
	Reader       string
	UserData     interface{}
	CurrentState StateFlag
	EventState   StateFlag
	// TODO: ATR
}
