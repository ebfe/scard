// Package scard provides bindings to the PC/SC API.
package scard

import (
	"unsafe"
)

type CardStatus struct {
	Reader         string
	State          State
	ActiveProtocol Protocol
	Atr            []byte
}

type ReaderState struct {
	Reader       string
	UserData     interface{}
	CurrentState StateFlag
	EventState   StateFlag
	Atr          []byte
}

func (buf strbuf) ptr() unsafe.Pointer {
	return unsafe.Pointer(&buf[0])
}

func (buf strbuf) split() []strbuf {
	var chunks []strbuf
	for len(buf) > 0 && buf[0] != 0 {
		i := 0
		for i = range buf {
			if buf[i] == 0 {
				break
			}
		}
		chunks = append(chunks, buf[:i+1])
		buf = buf[i+1:]
	}

	return chunks
}

func encodemstr(strings ...string) (strbuf, error) {
	var buf strbuf
	for _, s := range strings {
		utf16, err := encodestr(s)
		if err != nil {
			return nil, err
		}
		buf = append(buf, utf16...)
	}
	buf = append(buf, 0)
	return buf, nil
}

func decodemstr(buf strbuf) []string {
	var strings []string
	for _, chunk := range buf.split() {
		strings = append(strings, decodestr(chunk))
	}
	return strings
}
