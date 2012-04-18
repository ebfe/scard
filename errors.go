package scard

// #cgo pkg-config: libpcsclite
// #include <winscard.h>
import "C"

type scardError uint

func (e *scardError) Error() string {
	return "scard: " + C.GoString(C.pcsc_stringify_error(C.LONG(*e)))
}

func newError(code C.LONG) *scardError {
	e := scardError(code)
	return &e
}

const (
	S_SUCCESS                 scardError = C.SCARD_S_SUCCESS
	F_INTERNAL_ERROR          scardError = C.SCARD_F_INTERNAL_ERROR
	E_CANCELLED               scardError = C.SCARD_E_CANCELLED
	E_INVALID_HANDLE          scardError = C.SCARD_E_INVALID_HANDLE
	E_INVALID_PARAMETER       scardError = C.SCARD_E_INVALID_PARAMETER
	E_INVALID_TARGET          scardError = C.SCARD_E_INVALID_TARGET
	E_NO_MEMORY               scardError = C.SCARD_E_NO_MEMORY
	F_WAITED_TOO_LONG         scardError = C.SCARD_F_WAITED_TOO_LONG
	E_INSUFFICIENT_BUFFER     scardError = C.SCARD_E_INSUFFICIENT_BUFFER
	E_UNKNOWN_READER          scardError = C.SCARD_E_UNKNOWN_READER
	E_TIMEOUT                 scardError = C.SCARD_E_TIMEOUT
	E_SHARING_VIOLATION       scardError = C.SCARD_E_SHARING_VIOLATION
	E_NO_SMARTCARD            scardError = C.SCARD_E_NO_SMARTCARD
	E_UNKNOWN_CARD            scardError = C.SCARD_E_UNKNOWN_CARD
	E_CANT_DISPOSE            scardError = C.SCARD_E_CANT_DISPOSE
	E_PROTO_MISMATCH          scardError = C.SCARD_E_PROTO_MISMATCH
	E_NOT_READY               scardError = C.SCARD_E_NOT_READY
	E_INVALID_VALUE           scardError = C.SCARD_E_INVALID_VALUE
	E_SYSTEM_CANCELLED        scardError = C.SCARD_E_SYSTEM_CANCELLED
	F_COMM_ERROR              scardError = C.SCARD_F_COMM_ERROR
	F_UNKNOWN_ERROR           scardError = C.SCARD_F_UNKNOWN_ERROR
	E_INVALID_ATR             scardError = C.SCARD_E_INVALID_ATR
	E_NOT_TRANSACTED          scardError = C.SCARD_E_NOT_TRANSACTED
	E_READER_UNAVAILABLE      scardError = C.SCARD_E_READER_UNAVAILABLE
	P_SHUTDOWN                scardError = C.SCARD_P_SHUTDOWN
	E_PCI_TOO_SMALL           scardError = C.SCARD_E_PCI_TOO_SMALL
	E_READER_UNSUPPORTED      scardError = C.SCARD_E_READER_UNSUPPORTED
	E_DUPLICATE_READER        scardError = C.SCARD_E_DUPLICATE_READER
	E_CARD_UNSUPPORTED        scardError = C.SCARD_E_CARD_UNSUPPORTED
	E_NO_SERVICE              scardError = C.SCARD_E_NO_SERVICE
	E_SERVICE_STOPPED         scardError = C.SCARD_E_SERVICE_STOPPED
	E_UNEXPECTED              scardError = C.SCARD_E_UNEXPECTED
	E_UNSUPPORTED_FEATURE     scardError = C.SCARD_E_UNSUPPORTED_FEATURE
	E_ICC_INSTALLATION        scardError = C.SCARD_E_ICC_INSTALLATION
	E_ICC_CREATEORDER         scardError = C.SCARD_E_ICC_CREATEORDER
	E_DIR_NOT_FOUND           scardError = C.SCARD_E_DIR_NOT_FOUND
	E_FILE_NOT_FOUND          scardError = C.SCARD_E_FILE_NOT_FOUND
	E_NO_DIR                  scardError = C.SCARD_E_NO_DIR
	E_NO_FILE                 scardError = C.SCARD_E_NO_FILE
	E_NO_ACCESS               scardError = C.SCARD_E_NO_ACCESS
	E_WRITE_TOO_MANY          scardError = C.SCARD_E_WRITE_TOO_MANY
	E_BAD_SEEK                scardError = C.SCARD_E_BAD_SEEK
	E_INVALID_CHV             scardError = C.SCARD_E_INVALID_CHV
	E_UNKNOWN_RES_MNG         scardError = C.SCARD_E_UNKNOWN_RES_MNG
	E_NO_SUCH_CERTIFICATE     scardError = C.SCARD_E_NO_SUCH_CERTIFICATE
	E_CERTIFICATE_UNAVAILABLE scardError = C.SCARD_E_CERTIFICATE_UNAVAILABLE
	E_NO_READERS_AVAILABLE    scardError = C.SCARD_E_NO_READERS_AVAILABLE
	E_COMM_DATA_LOST          scardError = C.SCARD_E_COMM_DATA_LOST
	E_NO_KEY_CONTAINER        scardError = C.SCARD_E_NO_KEY_CONTAINER
	E_SERVER_TOO_BUSY         scardError = C.SCARD_E_SERVER_TOO_BUSY
	W_UNSUPPORTED_CARD        scardError = C.SCARD_W_UNSUPPORTED_CARD
	W_UNRESPONSIVE_CARD       scardError = C.SCARD_W_UNRESPONSIVE_CARD
	W_UNPOWERED_CARD          scardError = C.SCARD_W_UNPOWERED_CARD
	W_RESET_CARD              scardError = C.SCARD_W_RESET_CARD
	W_REMOVED_CARD            scardError = C.SCARD_W_REMOVED_CARD
	W_SECURITY_VIOLATION      scardError = C.SCARD_W_SECURITY_VIOLATION
	W_WRONG_CHV               scardError = C.SCARD_W_WRONG_CHV
	W_CHV_BLOCKED             scardError = C.SCARD_W_CHV_BLOCKED
	W_EOF                     scardError = C.SCARD_W_EOF
	W_CANCELLED_BY_USER       scardError = C.SCARD_W_CANCELLED_BY_USER
	W_CARD_NOT_AUTHENTICATED  scardError = C.SCARD_W_CARD_NOT_AUTHENTICATED
)
