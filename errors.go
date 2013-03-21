package scard

// #cgo pkg-config: libpcsclite
// #include <winscard.h>
import "C"

type scardError uint32

func (e *scardError) Error() string {
	return "scard: " + C.GoString(C.pcsc_stringify_error(C.LONG(*e)))
}

func newError(code C.LONG) *scardError {
	e := scardError(code)
	return &e
}

const (
	S_SUCCESS                 scardError = 0x00000000
	F_INTERNAL_ERROR          scardError = 0x80100001
	E_CANCELLED               scardError = 0x80100002
	E_INVALID_HANDLE          scardError = 0x80100003
	E_INVALID_PARAMETER       scardError = 0x80100004
	E_INVALID_TARGET          scardError = 0x80100005
	E_NO_MEMORY               scardError = 0x80100006
	F_WAITED_TOO_LONG         scardError = 0x80100007
	E_INSUFFICIENT_BUFFER     scardError = 0x80100008
	E_UNKNOWN_READER          scardError = 0x80100009
	E_TIMEOUT                 scardError = 0x8010000A
	E_SHARING_VIOLATION       scardError = 0x8010000B
	E_NO_SMARTCARD            scardError = 0x8010000C
	E_UNKNOWN_CARD            scardError = 0x8010000D
	E_CANT_DISPOSE            scardError = 0x8010000E
	E_PROTO_MISMATCH          scardError = 0x8010000F
	E_NOT_READY               scardError = 0x80100010
	E_INVALID_VALUE           scardError = 0x80100011
	E_SYSTEM_CANCELLED        scardError = 0x80100012
	F_COMM_ERROR              scardError = 0x80100013
	F_UNKNOWN_ERROR           scardError = 0x80100014
	E_INVALID_ATR             scardError = 0x80100015
	E_NOT_TRANSACTED          scardError = 0x80100016
	E_READER_UNAVAILABLE      scardError = 0x80100017
	P_SHUTDOWN                scardError = 0x80100018
	E_PCI_TOO_SMALL           scardError = 0x80100019
	E_READER_UNSUPPORTED      scardError = 0x8010001A
	E_DUPLICATE_READER        scardError = 0x8010001B
	E_CARD_UNSUPPORTED        scardError = 0x8010001C
	E_NO_SERVICE              scardError = 0x8010001D
	E_SERVICE_STOPPED         scardError = 0x8010001E
	E_UNEXPECTED              scardError = 0x8010001F
	E_UNSUPPORTED_FEATURE     scardError = 0x8010001F
	E_ICC_INSTALLATION        scardError = 0x80100020
	E_ICC_CREATEORDER         scardError = 0x80100021
	E_FILE_NOT_FOUND          scardError = 0x80100024
	E_NO_DIR                  scardError = 0x80100025
	E_NO_FILE                 scardError = 0x80100026
	E_NO_ACCESS               scardError = 0x80100027
	E_WRITE_TOO_MANY          scardError = 0x80100028
	E_BAD_SEEK                scardError = 0x80100029
	E_INVALID_CHV             scardError = 0x8010002A
	E_UNKNOWN_RES_MNG         scardError = 0x8010002B
	E_NO_SUCH_CERTIFICATE     scardError = 0x8010002C
	E_CERTIFICATE_UNAVAILABLE scardError = 0x8010002D
	E_NO_READERS_AVAILABLE    scardError = 0x8010002E
	E_COMM_DATA_LOST          scardError = 0x8010002F
	E_NO_KEY_CONTAINER        scardError = 0x80100030
	E_SERVER_TOO_BUSY         scardError = 0x80100031
	W_UNSUPPORTED_CARD        scardError = 0x80100065
	W_UNRESPONSIVE_CARD       scardError = 0x80100066
	W_UNPOWERED_CARD          scardError = 0x80100067
	W_RESET_CARD              scardError = 0x80100068
	W_REMOVED_CARD            scardError = 0x80100069
	W_SECURITY_VIOLATION      scardError = 0x8010006A
	W_WRONG_CHV               scardError = 0x8010006B
	W_CHV_BLOCKED             scardError = 0x8010006C
	W_EOF                     scardError = 0x8010006D
	W_CANCELLED_BY_USER       scardError = 0x8010006E
	W_CARD_NOT_AUTHENTICATED  scardError = 0x8010006F
)
