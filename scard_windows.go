package scard

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	modwinscard = syscall.NewLazyDLL("winscard.dll")

	procEstablishContext = modwinscard.NewProc("SCardEstablishContext")
	procReleaseContext   = modwinscard.NewProc("SCardReleaseContext")
	procIsValidContext   = modwinscard.NewProc("SCardIsValidContext")
	procCancel           = modwinscard.NewProc("SCardCancel")
	procListReaders      = modwinscard.NewProc("SCardListReadersW")
	procListReaderGroups = modwinscard.NewProc("SCardListReaderGroupsW")
	procGetStatusChange  = modwinscard.NewProc("SCardGetStatusChangeW")
	procConnect          = modwinscard.NewProc("SCardConnectW")
	procDisconnect       = modwinscard.NewProc("SCardDisconnect")
	procReconnect        = modwinscard.NewProc("SCardReconnect")
	procBeginTransaction = modwinscard.NewProc("SCardBeginTransaction")
	procEndTransaction   = modwinscard.NewProc("SCardEndTransaction")
	procStatus           = modwinscard.NewProc("SCardStatusW")
	procTransmit         = modwinscard.NewProc("SCardTransmit")
	procControl          = modwinscard.NewProc("SCardControl")
	procGetAttrib        = modwinscard.NewProc("SCardGetAttrib")
	procSetAttrib        = modwinscard.NewProc("SCardSetAttrib")

	dataT0Pci = modwinscard.NewProc("g_rgSCardT0Pci")
	dataT1Pci = modwinscard.NewProc("g_rgSCardT1Pci")
)

var scardIoReqT0 uintptr
var scardIoReqT1 uintptr

func init() {
	if err := dataT0Pci.Find(); err != nil {
		panic(err)
	}
	scardIoReqT0 = dataT0Pci.Addr()
	if err := dataT1Pci.Find(); err != nil {
		panic(err)
	}
	scardIoReqT1 = dataT1Pci.Addr()
}

func (e Error) Error() string {
	err := syscall.Errno(e)
	return fmt.Sprintf("scard: error(%x): %s", uintptr(e), err.Error())
}

func scardEstablishContext(scope uint32, reserved1, reserved2 uintptr) (uintptr, Error) {
	var ctx uintptr
	r, _, _ := procEstablishContext.Call(uintptr(scope), reserved1, reserved2, uintptr(unsafe.Pointer(&ctx)))
	return ctx, Error(r)
}

func scardIsValidContext(ctx uintptr) Error {
	r, _, _ := procIsValidContext.Call(ctx)
	return Error(r)
}

func scardCancel(ctx uintptr) Error {
	r, _, _ := procCancel.Call(ctx)
	return Error(r)
}

func scardReleaseContext(ctx uintptr) Error {
	r, _, _ := procReleaseContext.Call(ctx)
	return Error(r)
}

func scardListReaders(ctx uintptr, groups, buf unsafe.Pointer, bufLen uint32) (uint32, Error) {
	dwBufLen := uint32(bufLen)
	r, _, _ := procListReaders.Call(ctx, uintptr(groups), uintptr(buf), uintptr(unsafe.Pointer(&dwBufLen)))
	return dwBufLen, Error(r)
}

func scardListReaderGroups(ctx uintptr, buf unsafe.Pointer, bufLen uint32) (uint32, Error) {
	dwBufLen := uint32(bufLen)
	r, _, _ := procListReaderGroups.Call(ctx, uintptr(buf), uintptr(unsafe.Pointer(&dwBufLen)))
	return dwBufLen, Error(r)
}

func scardGetStatusChange(ctx uintptr, timeout uint32, states []scardReaderState) Error {
	r, _, _ := procGetStatusChange.Call(ctx, uintptr(timeout), uintptr(unsafe.Pointer(&states[0])), uintptr(len(states)))
	return Error(r)
}

func scardConnect(ctx uintptr, reader unsafe.Pointer, shareMode ShareMode, proto Protocol) (uintptr, Protocol, Error) {
	var handle uintptr
	var activeProto uint32

	r, _, _ := procConnect.Call(ctx, uintptr(reader), uintptr(shareMode), uintptr(proto), uintptr(unsafe.Pointer(&handle)), uintptr(unsafe.Pointer(&activeProto)))

	return handle, Protocol(activeProto), Error(r)
}

func scardDisconnect(card uintptr, d Disposition) Error {
	r, _, _ := procDisconnect.Call(card, uintptr(d))
	return Error(r)
}

func scardReconnect(card uintptr, mode ShareMode, proto Protocol, disp Disposition) (Protocol, Error) {
	var activeProtocol uint32
	r, _, _ := procReconnect.Call(card, uintptr(mode), uintptr(proto), uintptr(disp), uintptr(unsafe.Pointer(&activeProtocol)))
	return Protocol(activeProtocol), Error(r)
}

func scardBeginTransaction(card uintptr) Error {
	r, _, _ := procBeginTransaction.Call(card)
	return Error(r)
}

func scardEndTransaction(card uintptr, disp Disposition) Error {
	r, _, _ := procEndTransaction.Call(card, uintptr(disp))
	return Error(r)
}

func scardCardStatus(card uintptr) (string, State, Protocol, []byte, Error) {
	var state, proto uint32
	var atr [maxAtrSize]byte
	var atrLen = uint32(len(atr))

	reader := make(strbuf, maxReadername+1)
	readerLen := uint32(len(reader))

	r, _, _ := procStatus.Call(card, uintptr(reader.ptr()), uintptr(unsafe.Pointer(&readerLen)), uintptr(unsafe.Pointer(&state)), uintptr(unsafe.Pointer(&proto)), uintptr(unsafe.Pointer(&atr[0])), uintptr(unsafe.Pointer(&atrLen)))

	return decodestr(reader[:readerLen]), State(state), Protocol(proto), atr[:atrLen], Error(r)
}

func scardTransmit(card uintptr, proto Protocol, cmd []byte, rsp []byte) (uint32, Error) {
	var sendpci uintptr
	var rspLen = uint32(len(rsp))

	switch proto {
	case ProtocolT0:
		sendpci = scardIoReqT0
	case ProtocolT1:
		sendpci = scardIoReqT1
	default:
		panic("unknown protocol")
	}

	r, _, _ := procTransmit.Call(card, sendpci, uintptr(unsafe.Pointer(&cmd[0])), uintptr(len(cmd)), uintptr(0), uintptr(unsafe.Pointer(&rsp[0])), uintptr(unsafe.Pointer(&rspLen)))

	return rspLen, Error(r)
}

func scardControl(card uintptr, ioctl uint32, in, out []byte) (uint32, Error) {
	var ptrIn uintptr
	var outLen = uint32(len(out))

	if len(in) != 0 {
		ptrIn = uintptr(unsafe.Pointer(&in[0]))
	}

	r, _, _ := procControl.Call(card, uintptr(ioctl), ptrIn, uintptr(len(in)), uintptr(unsafe.Pointer(&out[0])), uintptr(len(out)), uintptr(unsafe.Pointer(&outLen)))
	return outLen, Error(r)
}

func scardGetAttrib(card uintptr, id Attrib, buf []byte) (uint32, Error) {
	var ptr uintptr

	if len(buf) != 0 {
		ptr = uintptr(unsafe.Pointer(&buf[0]))
	}

	bufLen := uint32(len(buf))
	r, _, _ := procGetAttrib.Call(card, uintptr(id), ptr, uintptr(unsafe.Pointer(&bufLen)))

	return bufLen, Error(r)
}

func scardSetAttrib(card uintptr, id Attrib, buf []byte) Error {
	r, _, _ := procSetAttrib.Call(card, uintptr(id), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return Error(r)
}
/*
// wraps SCardEstablishContext
func EstablishContext() (*Context, error) {
	var ctx Context

	r, _, _ := procEstablishContext.Call(2, uintptr(0), uintptr(0), uintptr(unsafe.Pointer(&ctx.ctx)))
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	return &ctx, nil
}

// wraps SCardIsValidContext
func (ctx *Context) IsValid() (bool, error) {
	r, _, _ := procIsValidContext.Call(ctx.ctx)
	if Error(r) != ErrSuccess {
		return false, Error(r)
	}
	return true, nil
}

// wraps SCardCancel
func (ctx *Context) Cancel() error {
	r, _, _ := procCancel.Call(ctx.ctx)
	if Error(r) != ErrSuccess {
		return Error(r)
	}
	return nil
}

// wraps SCardReleaseContext
func (ctx *Context) Release() error {
	r, _, _ := procReleaseContext.Call(uintptr(ctx.ctx))
	if Error(r) != ErrSuccess {
		return Error(r)
	}
	return nil
}

// wraps SCardListReaders
func (ctx *Context) ListReaders() ([]string, error) {
	var needed uintptr

	r, _, _ := procListReaders.Call(
		ctx.ctx,
		0,
		0,
		uintptr(unsafe.Pointer(&needed)))
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	data := make(strbuf, needed)
	r, _, _ = procListReaders.Call(
		ctx.ctx,
		0,
		uintptr(data.ptr()),
		uintptr(unsafe.Pointer(&needed)))
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	return decodemstr(data[:needed]), nil
}

// wraps SCardListReaderGroups
func (ctx *Context) ListReaderGroups() ([]string, error) {
	var needed uintptr

	r, _, _ := procListReaderGroups.Call(
		ctx.ctx,
		0,
		uintptr(unsafe.Pointer(&needed)))
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	data := make(strbuf, needed)
	r, _, _ = procListReaderGroups.Call(
		ctx.ctx,
		uintptr(data.ptr()),
		uintptr(unsafe.Pointer(&needed)))
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	return decodemstr(data[:needed]), nil
}

// wraps SCardGetStatusChange
func (ctx *Context) GetStatusChange(readerStates []ReaderState, timeout time.Duration) error {

	dwTimeout := durationToTimeout(timeout)
	crs := make([]scardReaderState, len(readerStates))

	for i := range readerStates {
		buf, err := encodestr(readerStates[i].Reader)
		if err != nil {
			return err
		}
		crs[i].szReader = uintptr(buf.ptr())
		crs[i].dwCurrentState = uint32(readerStates[i].CurrentState)
		crs[i].cbAtr = uint32(len(readerStates[i].Atr))
		copy(crs[i].rgbAtr[:], readerStates[i].Atr)
	}

	r, _, _ := procGetStatusChange.Call(
		ctx.ctx,
		uintptr(dwTimeout),
		uintptr(unsafe.Pointer(&crs[0])),
		uintptr(len(crs)))

	if Error(r) != ErrSuccess {
		return Error(r)
	}

	for i := range readerStates {
		readerStates[i].EventState = StateFlag(crs[i].dwEventState)
		if crs[i].cbAtr > 0 {
			readerStates[i].Atr = make([]byte, int(crs[i].cbAtr))
			copy(readerStates[i].Atr, crs[i].rgbAtr[:crs[i].cbAtr])
		}
	}

	return nil
}

// wraps SCardConnect
func (ctx *Context) Connect(reader string, mode ShareMode, proto Protocol) (*Card, error) {
	var handle uintptr
	var activeProtocol uintptr

	creader, err := encodemstr(reader)
	if err != nil {
		return nil, err
	}

	r, _, _ := procConnect.Call(
		ctx.ctx,
		uintptr(creader.ptr()),
		uintptr(mode),
		uintptr(proto),
		uintptr(unsafe.Pointer(handle)),
		uintptr(unsafe.Pointer(activeProtocol)))

	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	return &Card{handle: handle, activeProtocol: Protocol(activeProtocol)}, nil
}

// wraps SCardDisconnect
func (card *Card) Disconnect(d Disposition) error {
	r, _, _ := procDisconnect.Call(card.handle, uintptr(d))
	if Error(r) != ErrSuccess {
		return Error(r)
	}
	return nil
}

// wraps SCardReconnect
func (card *Card) Reconnect(mode ShareMode, protocol Protocol, init Disposition) error {
	r, _, _ := procReconnect.Call(card.handle, uintptr(protocol), uintptr(init))
	if Error(r) != ErrSuccess {
		return Error(r)
	}
	return nil
}

// wraps SCardBeginTransaction
func (card *Card) BeginTransaction() error {
	r, _, _ := procBeginTransaction.Call(card.handle)
	if Error(r) != ErrSuccess {
		return Error(r)
	}
	return nil
}

// wraps SCardEndTransaction
func (card *Card) EndTransaction(d Disposition) error {
	r, _, _ := procEndTransaction.Call(card.handle, uintptr(d))
	if Error(r) != ErrSuccess {
		return Error(r)
	}
	return nil
}

// wraps SCardStatus
func (card *Card) Status() (*CardStatus, error) {
	var state, proto uint32
	var atr [maxAtrSize]byte
	var atrLen = uint32(len(atr))

	reader := make(strbuf, maxReadername+1)
	readerLen := len(reader)

	r, _, _ := procStatus.Call(
		card.handle,
		uintptr(reader.ptr()),
		uintptr(unsafe.Pointer(&readerLen)),
		uintptr(unsafe.Pointer(&state)),
		uintptr(unsafe.Pointer(&proto)),
		uintptr(unsafe.Pointer(&atr[0])),
		uintptr(unsafe.Pointer(&atrLen)))
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	reader = reader[:readerLen]

	status := &CardStatus{
		Reader:         decodestr(reader),
		State:          State(state),
		ActiveProtocol: Protocol(proto),
		Atr:            atr[0:atrLen],
	}

	return status, nil
}

// wraps SCardTransmit
func (card *Card) Transmit(cmd []byte) ([]byte, error) {
	var sendpci uintptr

	switch card.activeProtocol {
	case ProtocolT0:
		sendpci = scardIoReqT0
	case ProtocolT1:
		sendpci = scardIoReqT1
	default:
		panic("unknown protocol")
	}

	var recv [maxBufferSizeExtended]byte
	var recvlen = uint32(len(recv))

	r, _, _ := procTransmit.Call(card.handle,
		sendpci,
		uintptr(unsafe.Pointer(&cmd[0])),
		uintptr(len(cmd)),
		0,
		uintptr(unsafe.Pointer(&recv[0])),
		uintptr(unsafe.Pointer(&recvlen)))
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
	var recvlen uintptr
	var r uintptr

	if len(cmd) == 0 {
		r, _, _ = procControl.Call(
			card.handle,
			uintptr(ctrl),
			0,
			0,
			uintptr(unsafe.Pointer(&recv[0])),
			uintptr(len(recv)),
			uintptr(unsafe.Pointer(&recvlen)))
	} else {
		r, _, _ = procControl.Call(
			card.handle,
			uintptr(ctrl),
			uintptr(unsafe.Pointer(&cmd[0])),
			uintptr(len(cmd)),
			uintptr(unsafe.Pointer(&recv[0])),
			uintptr(len(recv)),
			uintptr(unsafe.Pointer(&recvlen)))
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
	var needed uintptr

	r, _, _ := procGetAttrib.Call(
		card.handle,
		uintptr(id),
		0,
		uintptr(unsafe.Pointer(&needed)))
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	var attrib = make([]byte, needed)

	r, _, _ = procGetAttrib.Call(
		card.handle,
		uintptr(id),
		uintptr(unsafe.Pointer(&attrib[0])),
		uintptr(unsafe.Pointer(&needed)))
	if Error(r) != ErrSuccess {
		return nil, Error(r)
	}

	return attrib[0:needed], nil
}

// wraps SCardSetAttrib
func (card *Card) SetAttrib(id Attrib, data []byte) error {
	r, _, _ := procSetAttrib.Call(
		card.handle,
		uintptr(id),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)))
	if Error(r) != ErrSuccess {
		return Error(r)
	}
	return nil
}
*/

type scardReaderState struct {
	szReader       uintptr
	pvUserData     uintptr
	dwCurrentState uint32
	dwEventState   uint32
	cbAtr          uint32
	rgbAtr         [36]byte
}

func (rs *ReaderState) toSys() (scardReaderState, error) {
	var sys scardReaderState
	creader, err := encodestr(rs.Reader)
	if err != nil {
		return scardReaderState{}, err
	}
	sys.szReader = uintptr(creader.ptr())
	sys.dwCurrentState = uint32(rs.CurrentState)
	sys.cbAtr = uint32(len(rs.Atr))
	copy(sys.rgbAtr[:], rs.Atr)
	return sys, nil
}

func (rs *ReaderState) update(sys *scardReaderState) {
	rs.EventState = StateFlag(sys.dwEventState)
	if sys.cbAtr > 0 {
		rs.Atr = make([]byte, int(sys.cbAtr))
		copy(rs.Atr, sys.rgbAtr[:])
	}
}

type strbuf []uint16

func encodestr(s string) (strbuf, error) {
	utf16, err := syscall.UTF16FromString(s)
	return strbuf(utf16), err
}

func decodestr(buf strbuf) string {
	return syscall.UTF16ToString(buf)
}
