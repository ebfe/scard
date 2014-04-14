package scard

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
