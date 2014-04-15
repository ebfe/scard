// +build !windows

package scard

// BUG(mg): Does not work on darwin. (older/different libpcsclite?)

// #cgo pkg-config: libpcsclite
// #include <stdlib.h>
// #include <winscard.h>
import "C"

import (
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

func scardEstablishContext(scope uint32, reserved1, reserved2 uintptr) (uintptr, Error) {
	var ctx C.SCARDCONTEXT
	r := C.SCardEstablishContext(C.DWORD(scope), C.LPCVOID(reserved1), C.LPCVOID(reserved2), &ctx)
	return uintptr(ctx), Error(r)
}

func scardIsValidContext(ctx uintptr) Error {
	r := C.SCardIsValidContext(C.SCARDCONTEXT(ctx))
	return Error(r)
}

func scardCancel(ctx uintptr) Error {
	r := C.SCardCancel(C.SCARDCONTEXT(ctx))
	return Error(r)
}

func scardReleaseContext(ctx uintptr) Error {
	r := C.SCardReleaseContext(C.SCARDCONTEXT(ctx))
	return Error(r)
}

func scardListReaders(ctx uintptr, groups, buf unsafe.Pointer, bufLen uint32) (uint32, Error) {
	dwBufLen := C.DWORD(bufLen)
	r := C.SCardListReaders(C.SCARDCONTEXT(ctx), (C.LPCSTR)(groups), (C.LPSTR)(buf), &dwBufLen)
	return uint32(dwBufLen), Error(r)
}

func scardListReaderGroups(ctx uintptr, buf unsafe.Pointer, bufLen uint32) (uint32, Error) {
	dwBufLen := C.DWORD(bufLen)
	r := C.SCardListReaderGroups(C.SCARDCONTEXT(ctx), (C.LPSTR)(buf), &dwBufLen)
	return uint32(dwBufLen), Error(r)
}

type scardReaderState C.SCARD_READERSTATE

func scardGetStatusChange(ctx uintptr, timeout uint32, states []scardReaderState) Error {
	r := C.SCardGetStatusChange(C.SCARDCONTEXT(ctx), C.DWORD(timeout), (C.LPSCARD_READERSTATE)(unsafe.Pointer(&states[0])), C.DWORD(len(states)))
	return Error(r)
}

func scardConnect(ctx uintptr, reader unsafe.Pointer, shareMode ShareMode, proto Protocol) (uintptr, Protocol, Error) {
	var handle C.SCARDHANDLE
	var activeProto C.DWORD

	r := C.SCardConnect(C.SCARDCONTEXT(ctx), C.LPCSTR(reader), C.DWORD(shareMode), C.DWORD(proto), &handle, &activeProto)

	return uintptr(handle), Protocol(activeProto), Error(r)
}

func scardDisconnect(card uintptr, d Disposition) Error {
	r := C.SCardDisconnect(C.SCARDHANDLE(card), C.DWORD(d))
	return Error(r)
}

func scardReconnect(card uintptr, mode ShareMode, proto Protocol, disp Disposition) (Protocol, Error) {
	var activeProtocol C.DWORD
	r := C.SCardReconnect(C.SCARDHANDLE(card), C.DWORD(mode), C.DWORD(proto), C.DWORD(disp), &activeProtocol)
	return Protocol(activeProtocol), Error(r)
}

func scardBeginTransaction(card uintptr) Error {
	r := C.SCardBeginTransaction(C.SCARDHANDLE(card))
	return Error(r)
}

func scardEndTransaction(card uintptr, disp Disposition) Error {
	r := C.SCardEndTransaction(C.SCARDHANDLE(card), C.DWORD(disp))
	return Error(r)
}

func scardCardStatus(card uintptr) (string, State, Protocol, []byte, Error) {
	var readerBuf [C.MAX_READERNAME + 1]byte
	var readerLen = C.DWORD(len(readerBuf))
	var state, proto C.DWORD
	var atr [maxAtrSize]byte
	var atrLen = C.DWORD(len(atr))

	r := C.SCardStatus(C.SCARDHANDLE(card), (C.LPSTR)(unsafe.Pointer(&readerBuf[0])), &readerLen, &state, &proto, (*C.BYTE)(&atr[0]), &atrLen)

	return decodestr(readerBuf[:readerLen]), State(state), Protocol(proto), atr[:atrLen], Error(r)
}

func scardTransmit(card uintptr, proto Protocol, cmd []byte, rsp []byte) (uint32, Error) {
	var sendpci C.SCARD_IO_REQUEST
	var recvpci C.SCARD_IO_REQUEST
	var rspLen = C.DWORD(len(rsp))

	switch proto {
	case ProtocolT0, ProtocolT1:
		sendpci.dwProtocol = C.ulong(proto)
	default:
		panic("unknown protocol")
	}
	sendpci.cbPciLength = C.sizeof_SCARD_IO_REQUEST

	r := C.SCardTransmit(C.SCARDHANDLE(card), &sendpci, (*C.BYTE)(&cmd[0]), C.DWORD(len(cmd)), &recvpci, (*C.BYTE)(&rsp[0]), &rspLen)

	return uint32(rspLen), Error(r)
}

func scardControl(card uintptr, ioctl uint32, in, out []byte) (uint32, Error) {
	var ptrIn C.LPCVOID
	var outLen = C.DWORD(len(out))

	if len(in) != 0 {
		ptrIn = C.LPCVOID(unsafe.Pointer(&in[0]))
	}

	r := C.SCardControl(C.SCARDHANDLE(card), C.DWORD(ioctl), ptrIn, C.DWORD(len(in)), (C.LPVOID)(unsafe.Pointer(&out[0])), C.DWORD(len(out)), &outLen)
	return uint32(outLen), Error(r)
}

func scardGetAttrib(card uintptr, id Attrib, buf []byte) (uint32, Error) {
	var ptr C.LPBYTE

	if len(buf) != 0 {
		ptr = C.LPBYTE(unsafe.Pointer(&buf[0]))
	}

	bufLen := C.DWORD(len(buf))
	r := C.SCardGetAttrib(C.SCARDHANDLE(card), C.DWORD(id), ptr, &bufLen)

	return uint32(bufLen), Error(r)
}

func scardSetAttrib(card uintptr, id Attrib, buf []byte) Error {
	r := C.SCardSetAttrib(C.SCARDHANDLE(card), C.DWORD(id), (C.LPBYTE)(unsafe.Pointer(&buf[0])), C.DWORD(len(buf)))
	return Error(r)
}

// wraps SCardEstablishContext
func EstablishContext() (*Context, error) {
	ctx, r := scardEstablishContext(C.SCARD_SCOPE_SYSTEM, 0, 0)
	if r != ErrSuccess {
		return nil, r
	}

	return &Context{ctx: ctx}, nil
}

// wraps SCardIsValidContext
func (ctx *Context) IsValid() (bool, error) {
	r := scardIsValidContext(ctx.ctx)
	switch r {
	case ErrSuccess:
		return true, nil
	case ErrInvalidHandle:
		return false, nil
	default:
		return false, r
	}
	panic("unreachable")
}

// wraps SCardCancel
func (ctx *Context) Cancel() error {
	r := scardCancel(ctx.ctx)
	if r != ErrSuccess {
		return r
	}
	return nil
}

// wraps SCardReleaseContext
func (ctx *Context) Release() error {
	r := scardReleaseContext(ctx.ctx)
	if r != ErrSuccess {
		return r
	}
	return nil
}

// wraps SCardListReaders
func (ctx *Context) ListReaders() ([]string, error) {
	needed, r := scardListReaders(ctx.ctx, nil, nil, 0)
	if r != ErrSuccess {
		return nil, r
	}

	buf := make(strbuf, needed)
	n, r := scardListReaders(ctx.ctx, nil, buf.ptr(), uint32(len(buf)))
	if r != ErrSuccess {
		return nil, r
	}
	return decodemstr(buf[:n]), nil
}

// wraps SCardListReaderGroups
func (ctx *Context) ListReaderGroups() ([]string, error) {
	needed, r := scardListReaderGroups(ctx.ctx, nil, 0)
	if r != ErrSuccess {
		return nil, r
	}

	buf := make(strbuf, needed)
	n, r := scardListReaderGroups(ctx.ctx, buf.ptr(), uint32(len(buf)))
	if r != ErrSuccess {
		return nil, r
	}
	return decodemstr(buf[:n]), nil
}

// wraps SCardGetStatusChange
func (ctx *Context) GetStatusChange(readerStates []ReaderState, timeout time.Duration) error {

	dwTimeout := durationToTimeout(timeout)
	states := make([]scardReaderState, len(readerStates))

	for i := range readerStates {
		states[i].szReader = (*C.char)(strbuf(readerStates[i].Reader).ptr())
		states[i].dwCurrentState = C.DWORD(readerStates[i].CurrentState)
		states[i].cbAtr = C.DWORD(len(readerStates[i].Atr))
		for j, b := range readerStates[i].Atr {
			states[i].rgbAtr[j] = C.uchar(b)
		}
	}

	r := scardGetStatusChange(ctx.ctx, dwTimeout, states)
	if r != ErrSuccess {
		return r
	}

	for i := range readerStates {
		readerStates[i].EventState = StateFlag(states[i].dwEventState)
		if states[i].cbAtr > 0 {
			readerStates[i].Atr = make([]byte, int(states[i].cbAtr))
			for j := C.DWORD(0); j < states[i].cbAtr; j++ {
				readerStates[i].Atr[j] = byte(states[i].rgbAtr[j])
			}
		}
	}

	return nil
}

// wraps SCardConnect
func (ctx *Context) Connect(reader string, mode ShareMode, proto Protocol) (*Card, error) {
	creader := strbuf(reader).ptr()
	handle, activeProtocol, r := scardConnect(ctx.ctx, creader, mode, proto)
	if r != ErrSuccess {
		return nil, r
	}
	return &Card{handle: handle, activeProtocol: activeProtocol}, nil
}

// wraps SCardDisconnect
func (card *Card) Disconnect(d Disposition) error {
	r := scardDisconnect(card.handle, d)
	if r != ErrSuccess {
		return r
	}
	return nil
}

// wraps SCardReconnect
func (card *Card) Reconnect(mode ShareMode, proto Protocol, disp Disposition) error {
	activeProtocol, r := scardReconnect(card.handle, mode, proto, disp)
	if r != ErrSuccess {
		return r
	}
	card.activeProtocol = activeProtocol
	return nil
}

// wraps SCardBeginTransaction
func (card *Card) BeginTransaction() error {
	r := scardBeginTransaction(card.handle)
	if r != ErrSuccess {
		return r
	}
	return nil
}

// wraps SCardEndTransaction
func (card *Card) EndTransaction(disp Disposition) error {
	r := scardEndTransaction(card.handle, disp)
	if r != ErrSuccess {
		return r
	}
	return nil
}

// wraps SCardStatus
func (card *Card) Status() (*CardStatus, error) {
	reader, state, proto, atr, err := scardCardStatus(card.handle)
	if err != ErrSuccess {
		return nil, err
	}
	return &CardStatus{Reader: reader, State: state, ActiveProtocol: proto, Atr: atr}, nil
}

// wraps SCardTransmit
func (card *Card) Transmit(cmd []byte) ([]byte, error) {
	rsp := make([]byte, maxBufferSizeExtended)
	rspLen, err := scardTransmit(card.handle, card.activeProtocol, cmd, rsp)
	if err != ErrSuccess {
		return nil, err
	}
	return rsp[:rspLen], nil
}

// wraps SCardControl
func (card *Card) Control(ioctl uint32, in []byte) ([]byte, error) {
	var out [0xffff]byte
	outLen, err := scardControl(card.handle, ioctl, in, out[:])
	if err != ErrSuccess {
		return nil, err
	}
	return out[:outLen], nil
}

// wraps SCardGetAttrib
func (card *Card) GetAttrib(id Attrib) ([]byte, error) {
	needed, err := scardGetAttrib(card.handle, id, nil)
	if err != ErrSuccess {
		return nil, err
	}

	var attrib = make([]byte, needed)
	n, err := scardGetAttrib(card.handle, id, attrib)
	if err != ErrSuccess {
		return nil, err
	}
	return attrib[:n], nil
}

// wraps SCardSetAttrib
func (card *Card) SetAttrib(id Attrib, data []byte) error {
	err := scardSetAttrib(card.handle, id, data)
	if err != ErrSuccess {
		return err
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
