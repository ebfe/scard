// +build !windows

package scard

// BUG(mg): Does not work on darwin. (older/different libpcsclite?)

// #cgo pkg-config: libpcsclite
// #include <stdlib.h>
// #include <winscard.h>
import "C"

import (
	"bytes"
	"time"
	"unsafe"
)

func (e Error) Error() string {
	return "scard: " + C.GoString(C.pcsc_stringify_error(C.LONG(e)))
}

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

// wraps SCardEstablishContext
func EstablishContext() (*Context, error) {
	var ctx Context

	r := C.SCardEstablishContext(C.SCARD_SCOPE_SYSTEM, nil, nil, &ctx.ctx)
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	return &ctx, nil
}

// wraps SCardIsValidContext
func (ctx *Context) IsValid() (bool, error) {
	r := C.SCardIsValidContext(ctx.ctx)
	switch Error(r) {
	case ErrSuccess:
		return true, nil
	case ErrInvalidHandle:
		return false, nil
	default:
		return false, Error(r)
	}
	panic("unreachable")
}

// wraps SCardCancel
func (ctx *Context) Cancel() error {
	r := C.SCardCancel(ctx.ctx)
	if Error(r) != ErrSuccess {
		return Error(r)
	}
	return nil
}

// wraps SCardReleaseContext
func (ctx *Context) Release() error {
	r := C.SCardReleaseContext(ctx.ctx)
	if Error(r) != ErrSuccess {
		return Error(r)
	}
	return nil
}


// wraps SCardListReaders
func (ctx *Context) ListReaders() ([]string, error) {
	var needed C.DWORD
	r := C.SCardListReaders(ctx.ctx, nil, nil, &needed)
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	data := make(strbuf, needed)
	r = C.SCardListReaders(ctx.ctx, nil, C.LPSTR(data.ptr()), &needed)
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	return decodemstr(data), nil
}

// wraps SCardListReaderGroups
func (ctx *Context) ListReaderGroups() ([]string, error) {
	var needed C.DWORD

	r := C.SCardListReaderGroups(ctx.ctx, nil, &needed)
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	data := make(strbuf, needed)

	r = C.SCardListReaderGroups(ctx.ctx, C.LPSTR(data.ptr()), &needed)
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	return decodemstr(data), nil
}

// wraps SCardGetStatusChange
func (ctx *Context) GetStatusChange(readerStates []ReaderState, timeout time.Duration) error {

	var dwTimeout uint32

	switch {
	case timeout < 0:
		dwTimeout = infiniteTimeout
	case timeout > time.Duration(infiniteTimeout)*time.Millisecond:
		dwTimeout = infiniteTimeout - 1
	default:
		dwTimeout = uint32(timeout / time.Millisecond)
	}

	crs := make([]C.SCARD_READERSTATE, len(readerStates))

	for i := range readerStates {
		crs[i].szReader = (*C.char)(strbuf(readerStates[i].Reader).ptr())
		crs[i].dwCurrentState = C.DWORD(readerStates[i].CurrentState)
		crs[i].cbAtr = C.DWORD(len(readerStates[i].Atr))
		for j, b := range readerStates[i].Atr {
			crs[i].rgbAtr[j] = C.uchar(b)
		}
	}

	r := C.SCardGetStatusChange(ctx.ctx, C.DWORD(dwTimeout),
		(C.LPSCARD_READERSTATE)(unsafe.Pointer(&crs[0])),
		C.DWORD(len(crs)))

	if Error(r) != ErrSuccess {
		return Error(r)
	}

	for i := range readerStates {
		readerStates[i].EventState = StateFlag(crs[i].dwEventState)
		if crs[i].cbAtr > 0 {
			readerStates[i].Atr = make([]byte, int(crs[i].cbAtr))
			for j := C.DWORD(0); j < crs[i].cbAtr; j++ {
				readerStates[i].Atr[j] = byte(crs[i].rgbAtr[j])
			}
		}
	}

	return nil
}

// wraps SCardConnect
func (ctx *Context) Connect(reader string, mode ShareMode, proto Protocol) (*Card, error) {
	var card Card
	var activeProtocol C.DWORD

	creader := (*C.char)(strbuf(reader).ptr())

	r := C.SCardConnect(ctx.ctx, creader, C.DWORD(mode), C.DWORD(proto), &card.handle, &activeProtocol)
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	card.activeProtocol = Protocol(activeProtocol)
	return &card, nil
}

// wraps SCardDisconnect
func (card *Card) Disconnect(d Disposition) error {
	r := C.SCardDisconnect(card.handle, C.DWORD(d))
	if Error(r) != ErrSuccess {
		return Error(r)
	}
	return nil
}

// wraps SCardReconnect
func (card *Card) Reconnect(mode ShareMode, protocol Protocol, init Disposition) error {
	var activeProtocol C.DWORD

	r := C.SCardReconnect(card.handle, C.DWORD(mode), C.DWORD(protocol), C.DWORD(init), &activeProtocol)
	if Error(r) != ErrSuccess {
		return Error(r)
	}

	card.activeProtocol = Protocol(activeProtocol)

	return nil
}

// wraps SCardBeginTransaction
func (card *Card) BeginTransaction() error {
	r := C.SCardBeginTransaction(card.handle)
	if Error(r) != ErrSuccess {
		return Error(r)
	}
	return nil
}

// wraps SCardEndTransaction
func (card *Card) EndTransaction(d Disposition) error {
	r := C.SCardEndTransaction(card.handle, C.DWORD(d))
	if Error(r) != ErrSuccess {
		return Error(r)
	}
	return nil
}

// wraps SCardStatus
func (card *Card) Status() (*CardStatus, error) {
	var readerBuf [C.MAX_READERNAME + 1]byte
	var readerLen = C.DWORD(len(readerBuf))
	var state, proto C.DWORD
	var atr [maxAtrSize]byte
	var atrLen = C.DWORD(len(atr))

	r := C.SCardStatus(card.handle, (C.LPSTR)(unsafe.Pointer(&readerBuf[0])), &readerLen, &state, &proto, (*C.BYTE)(&atr[0]), &atrLen)
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	// strip terminating 0
	reader := readerBuf[:readerLen]
	if z := bytes.IndexByte(reader, 0); z != -1 {
		reader = reader[:z]
	}

	status := &CardStatus{
		Reader:         string(reader),
		State:          State(state),
		ActiveProtocol: Protocol(proto),
		Atr:            atr[0:atrLen],
	}

	return status, nil
}

// wraps SCardTransmit
func (card *Card) Transmit(cmd []byte) ([]byte, error) {
	var sendpci C.SCARD_IO_REQUEST
	var recvpci C.SCARD_IO_REQUEST

	switch card.activeProtocol {
	case ProtocolT0, ProtocolT1:
		sendpci.dwProtocol = C.ulong(card.activeProtocol)
	default:
		panic("unknown protocol")
	}
	sendpci.cbPciLength = C.sizeof_SCARD_IO_REQUEST

	var recv [maxBufferSizeExtended]byte
	var recvlen = C.DWORD(len(recv))

	r := C.SCardTransmit(card.handle, &sendpci, (*C.BYTE)(&cmd[0]), C.DWORD(len(cmd)), &recvpci, (*C.BYTE)(&recv[0]), &recvlen)
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	rsp := make([]byte, recvlen)
	copy(rsp, recv[0:recvlen])

	return rsp, nil
}

// wraps SCardControl
func (card *Card) Control(ctrl uint32, cmd []byte) ([]byte, error) {
	var recv [0xffff]byte
	var recvlen C.DWORD
	var r C.LONG

	if len(cmd) == 0 {
		r = C.SCardControl(card.handle, C.DWORD(ctrl),
			(C.LPCVOID)(nil), 0,
			(C.LPVOID)(unsafe.Pointer(&recv[0])), C.DWORD(len(recv)), &recvlen)
	} else {
		r = C.SCardControl(card.handle, C.DWORD(ctrl),
			(C.LPCVOID)(unsafe.Pointer(&cmd[0])), C.DWORD(len(cmd)),
			(C.LPVOID)(unsafe.Pointer(&recv[0])), C.DWORD(len(recv)), &recvlen)
	}
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	rsp := make([]byte, recvlen)
	copy(rsp, recv[0:recvlen])

	return rsp, nil
}

// wraps SCardGetAttrib
func (card *Card) GetAttrib(id Attrib) ([]byte, error) {
	var needed C.DWORD

	r := C.SCardGetAttrib(card.handle, C.DWORD(id), nil, &needed)
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	var attrib = make([]byte, needed)

	r = C.SCardGetAttrib(card.handle, C.DWORD(id), (*C.BYTE)(&attrib[0]), &needed)
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	return attrib[0:needed], nil
}

// wraps SCardSetAttrib
func (card *Card) SetAttrib(id Attrib, data []byte) error {
	r := C.SCardSetAttrib(card.handle, C.DWORD(id), (*C.BYTE)(&data[0]), C.DWORD(len(data)))
	if Error(r) != ErrSuccess {
		return Error(r)
	}
	return nil
}

type strbuf []byte

func encodestr(s string) (strbuf, error) {
	buf := strbuf(s + "\x00")
	return buf, nil
}

func decodestr(buf strbuf) string {
	if len(buf) == 0 {
		return ""
	}

	if buf[len(buf)-1] == 0 {
		buf = buf[:len(buf)-1]
	}

	return string(buf)
}
