package scard

type Protocol uint32

const (
	PROTOCOL_UNDEFINED Protocol = iota
	PROTOCOL_T0
	PROTOCOL_T1
	PROTOCOL_ANY = PROTOCOL_T0 | PROTOCOL_T1
)

type ShareMode uint32

const (
	SHARE_EXCLUSIVE ShareMode = iota + 1
	SHARE_SHARED
	SHARE_DIRECT
)

type Disposition uint32

const (
	LEAVE_CARD Disposition = iota
	RESET_CARD
	UNPOWER_CARD
	EJECT_CARD
)

type State uint16

const (
	UNKNOWN State = 0
	ABSENT  State = 1 << iota
	PRESENT
	SWALLOWED
	POWERED
	NEGOTIABLE
	SPECIFIC
)

func (s State) String() string {
	switch {
	case s&ABSENT == ABSENT:
		return "ABSENT"
	case s&PRESENT == PRESENT:
		return "PRESENT"
	case s&SWALLOWED == SWALLOWED:
		return "SWALLOWED"
	case s&POWERED == POWERED:
		return "POWERED"
	case s&NEGOTIABLE == NEGOTIABLE:
		return "NEGOTIABLE"
	case s&SPECIFIC == SPECIFIC:
		return "SPECIFIC"
	default:
		return "UNKNOWN"
	}
}

type StateFlag uint32

const (
	STATE_UNAWARE StateFlag = 0
	STATE_IGNORE  StateFlag = (1 << iota)
	STATE_CHANGED
	STATE_UNKNOWN
	STATE_UNAVAILABLE
	STATE_EMPTY
	STATE_PRESENT
	STATE_ATRMATCH
	STATE_EXCLUSIVE
	STATE_INUSE
	STATE_MUTE
	STATE_UNPOWERED
)

type Timeout uint32

const (
	INFINITE Timeout = 0xffffffff
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
