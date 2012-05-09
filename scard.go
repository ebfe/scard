// pcsclite binding
// 
// http://pcsclite.alioth.debian.org/pcsclite.html
// http://pcsclite.alioth.debian.org/api/group__API.html
//
package scard

// #cgo pkg-config: libpcsclite
// #include <stdlib.h>
// #include <winscard.h>
import "C"

import (
	"bytes"
	"unsafe"
)

// Version returns the libpcsclite version string
func Version() string {
	return C.PCSCLITE_VERSION_NUMBER
}

type Context struct {
	ctx C.SCARDCONTEXT
}

type Card struct {
	handle         C.SCARDHANDLE
	activeProtocol Protocol
}

type Protocol uint32

const (
	PROTOCOL_UNDEFINED Protocol = C.SCARD_PROTOCOL_UNDEFINED
	PROTOCOL_T0        Protocol = C.SCARD_PROTOCOL_T0
	PROTOCOL_T1        Protocol = C.SCARD_PROTOCOL_T1
	PROTOCOL_RAW       Protocol = C.SCARD_PROTOCOL_RAW
	PROTOCOL_ANY       Protocol = C.SCARD_PROTOCOL_ANY
)

type ShareMode uint32

const (
	SHARE_EXCLUSIVE ShareMode = C.SCARD_SHARE_EXCLUSIVE
	SHARE_SHARED    ShareMode = C.SCARD_SHARE_SHARED
	SHARE_DIRECT    ShareMode = C.SCARD_SHARE_DIRECT
)

type Disposition uint32

const (
	LEAVE_CARD   Disposition = C.SCARD_LEAVE_CARD
	RESET_CARD   Disposition = C.SCARD_RESET_CARD
	UNPOWER_CARD Disposition = C.SCARD_UNPOWER_CARD
	EJECT_CARD   Disposition = C.SCARD_EJECT_CARD
)

type CardStatus struct {
	Reader         string
	State          uint32
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

type StateFlag uint32

const (
	STATE_UNAWARE     StateFlag = C.SCARD_STATE_UNAWARE
	STATE_IGNORE      StateFlag = C.SCARD_STATE_IGNORE
	STATE_CHANGED     StateFlag = C.SCARD_STATE_CHANGED
	STATE_UNKNOWN     StateFlag = C.SCARD_STATE_UNKNOWN
	STATE_UNAVAILABLE StateFlag = C.SCARD_STATE_UNAVAILABLE
	STATE_EMPTY       StateFlag = C.SCARD_STATE_EMPTY
	STATE_PRESENT     StateFlag = C.SCARD_STATE_PRESENT
	STATE_ATRMATCH    StateFlag = C.SCARD_STATE_ATRMATCH
	STATE_EXCLUSIVE   StateFlag = C.SCARD_STATE_EXCLUSIVE
	STATE_INUSE       StateFlag = C.SCARD_STATE_INUSE
	STATE_MUTE        StateFlag = C.SCARD_STATE_MUTE
	STATE_UNPOWERED   StateFlag = C.SCARD_STATE_UNPOWERED
)

type Timeout uint32

const (
	INFINITE Timeout = C.INFINITE
)

// wraps SCardEstablishContext
func EstablishContext() (*Context, error) {
	var ctx Context

	r := C.SCardEstablishContext(C.SCARD_SCOPE_SYSTEM, nil, nil, &ctx.ctx)
	if r != C.SCARD_S_SUCCESS {
		return nil, newError(r)
	}

	return &ctx, nil
}

// wraps SCardIsValidContext
func (ctx *Context) IsValid() (bool, error) {
	r := C.SCardIsValidContext(ctx.ctx)
	switch r {
	case C.SCARD_S_SUCCESS:
		return true, nil
	case C.SCARD_E_INVALID_HANDLE:
		return false, nil
	default:
		return false, newError(r)
	}
	panic("unreachable")
}

// wraps SCardCancel
func (ctx *Context) Cancel() error {
	r := C.SCardCancel(ctx.ctx)
	if r != C.SCARD_S_SUCCESS {
		return newError(r)
	}
	return nil
}

// wraps SCardReleaseContext
func (ctx *Context) Release() error {
	r := C.SCardReleaseContext(ctx.ctx)
	if r != C.SCARD_S_SUCCESS {
		return newError(r)
	}
	return nil
}

// wraps SCardListReaders
func (ctx *Context) ListReaders() ([]string, error) {
	var needed C.DWORD

	r := C.SCardListReaders(ctx.ctx, nil, nil, &needed)
	if r != C.SCARD_S_SUCCESS {
		return nil, newError(r)
	}

	data := make([]byte, needed)
	cdata := (*C.char)(unsafe.Pointer(&data[0]))

	r = C.SCardListReaders(ctx.ctx, nil, cdata, &needed)
	if r != C.SCARD_S_SUCCESS {
		return nil, newError(r)
	}

	var readers []string
	for _, b := range bytes.Split(data, []byte{0}) {
		if len(b) > 0 {
			readers = append(readers, string(b))
		}
	}

	return readers, nil
}

// wraps SCardListReaderGroups
func (ctx *Context) ListReaderGroups() ([]string, error) {
	var needed C.DWORD

	r := C.SCardListReaderGroups(ctx.ctx, nil, &needed)
	if r != C.SCARD_S_SUCCESS {
		return nil, newError(r)
	}

	data := make([]byte, needed)
	cdata := (*C.char)(unsafe.Pointer(&data[0]))

	r = C.SCardListReaderGroups(ctx.ctx, cdata, &needed)
	if r != C.SCARD_S_SUCCESS {
		return nil, newError(r)
	}

	var groups []string
	for _, b := range bytes.Split(data, []byte{0}) {
		if len(b) > 0 {
			groups = append(groups, string(b))
		}
	}

	return groups, nil
}

// wraps SCardGetStatusChange
func (ctx *Context) GetStatusChange(readerStates []ReaderState, timeout Timeout) error {

	crs := make([]C.SCARD_READERSTATE, len(readerStates))

	for i := range readerStates {
		crs[i].szReader = C.CString(readerStates[i].Reader)
		defer C.free(unsafe.Pointer(crs[i].szReader))
		crs[i].dwCurrentState = C.DWORD(readerStates[i].CurrentState)
	}

	r := C.SCardGetStatusChange(ctx.ctx, C.DWORD(timeout),
		(C.LPSCARD_READERSTATE)(unsafe.Pointer(&crs[0])),
		C.DWORD(len(crs)))

	if r != C.SCARD_S_SUCCESS {
		return newError(r)
	}

	for i := range readerStates {
		readerStates[i].EventState = StateFlag(crs[i].dwEventState)
	}

	return nil
}

// wraps SCardConnect
func (ctx *Context) Connect(reader string, mode ShareMode, proto Protocol) (*Card, error) {
	var card Card
	var activeProtocol C.DWORD

	creader := C.CString(reader)
	defer C.free(unsafe.Pointer(creader))

	r := C.SCardConnect(ctx.ctx, creader, C.DWORD(mode), C.DWORD(proto), &card.handle, &activeProtocol)
	if r != C.SCARD_S_SUCCESS {
		return nil, newError(r)
	}

	card.activeProtocol = Protocol(activeProtocol)
	return &card, nil
}

// wraps SCardDisconnect
func (card *Card) Disconnect(d Disposition) error {
	r := C.SCardDisconnect(card.handle, C.DWORD(d))
	if r != C.SCARD_S_SUCCESS {
		return newError(r)
	}
	return nil
}

// wraps SCardReconnect
func (card *Card) Reconnect(mode ShareMode, protocol Protocol, init Disposition) error {
	var activeProtocol C.DWORD

	r := C.SCardReconnect(card.handle, C.DWORD(mode), C.DWORD(protocol), C.DWORD(init), &activeProtocol)
	if r != C.SCARD_S_SUCCESS {
		return newError(r)
	}

	card.activeProtocol = Protocol(activeProtocol)

	return nil
}

// wraps SCardBeginTransaction
func (card *Card) BeginTransaction() error {
	r := C.SCardBeginTransaction(card.handle)
	if r != C.SCARD_S_SUCCESS {
		return newError(r)
	}
	return nil
}

// wraps SCardEndTransaction
func (card *Card) EndTransaction(d Disposition) error {
	r := C.SCardEndTransaction(card.handle, C.DWORD(d))
	if r != C.SCARD_S_SUCCESS {
		return newError(r)
	}
	return nil
}

// wraps SCardStatus
func (card *Card) Status() (*CardStatus, error) {
	var reader [C.MAX_READERNAME + 1]byte
	var readerLen = C.DWORD(len(reader))
	var state, proto C.DWORD
	var atr [C.MAX_ATR_SIZE]byte
	var atrLen = C.DWORD(len(atr))

	r := C.SCardStatus(card.handle, (C.LPSTR)(unsafe.Pointer(&reader[0])), &readerLen, &state, &proto, (*C.BYTE)(&atr[0]), &atrLen)
	if r != C.SCARD_S_SUCCESS {
		return nil, newError(r)
	}

	status := &CardStatus{
		Reader:         string(reader[0:readerLen]),
		State:          uint32(state),
		ActiveProtocol: Protocol(proto),
		ATR:            atr[0:atrLen],
	}

	return status, nil
}

// wraps SCardTransmit
func (card *Card) Transmit(cmd []byte) ([]byte, error) {
	var sendpci *C.SCARD_IO_REQUEST
	var recvpci C.SCARD_IO_REQUEST

	switch card.activeProtocol {
	case PROTOCOL_T0:
		sendpci = &C.g_rgSCardT0Pci
	case PROTOCOL_T1:
		sendpci = &C.g_rgSCardT1Pci
	case PROTOCOL_RAW:
		sendpci = &C.g_rgSCardRawPci
	default:
		panic("unknown protocol")
	}

	var recv [C.MAX_BUFFER_SIZE_EXTENDED]byte
	var recvlen C.DWORD = C.DWORD(len(recv))

	r := C.SCardTransmit(card.handle, sendpci, (*C.BYTE)(&cmd[0]), C.DWORD(len(cmd)), &recvpci, (*C.BYTE)(&recv[0]), &recvlen)
	if r != C.SCARD_S_SUCCESS {
		return nil, newError(r)
	}

	rsp := make([]byte, recvlen)
	copy(rsp, recv[0:recvlen])

	return rsp, nil
}

// wraps SCardControl
func (card *Card) Control(ctrl uint32, cmd []byte) ([]byte, error) {
	var recv [C.MAX_BUFFER_SIZE_EXTENDED]byte
	var recvlen C.DWORD

	r := C.SCardControl(card.handle, C.DWORD(ctrl),
		(C.LPCVOID)(unsafe.Pointer(&cmd[0])), C.DWORD(len(cmd)),
		(C.LPVOID)(unsafe.Pointer(&recv[0])), C.DWORD(len(recv)), &recvlen)
	if r != C.SCARD_S_SUCCESS {
		return nil, newError(r)
	}

	rsp := make([]byte, recvlen)
	copy(rsp, recv[0:recvlen])

	return rsp, nil
}

// wraps SCardGetAttrib
func (card *Card) GetAttrib(id uint32) ([]byte, error) {
	var needed C.DWORD

	r := C.SCardGetAttrib(card.handle, C.DWORD(id), nil, &needed)
	if r != C.SCARD_S_SUCCESS {
		return nil, newError(r)
	}

	var attrib = make([]byte, needed)

	r = C.SCardGetAttrib(card.handle, C.DWORD(id), (*C.BYTE)(&attrib[0]), &needed)
	if r != C.SCARD_S_SUCCESS {
		return nil, newError(r)
	}

	return attrib[0:needed], nil
}

// wraps SCardSetAttrib
func (card *Card) SetAttrib(id uint32, data []byte) error {
	r := C.SCardSetAttrib(card.handle, C.DWORD(id), (*C.BYTE)(&data[0]), C.DWORD(len(data)))
	if r != C.SCARD_S_SUCCESS {
		return newError(r)
	}
	return nil
}

// SCardFreeMemory is not needed. We (hopefuly) never return buffers allocated by libpcsclite
