package utfutil

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
)

// Encode encodes UTF-16 to UTF-8
func Encode(dst io.Writer, src io.Reader, endines binary.ByteOrder) error {
	// TODO: remove full read, instead read chunks
	p, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}

	// handle BOM
	if p[0] == 0xfe && p[1] == 0xff {
		//endines = binary.BigEndian
		p = p[2:]
		//if _, err = dst.Write([]byte{0xef, 0xbb, 0xbf}); err != nil {
		//	return err
		//}
	}
	// else if p[0] == 0xff && p[1] == 0xfe {
	//endines = binary.LittleEndian
	//}

	var (
		i, j, n = 0, 0, len(p)
		out     = make([]byte, n*2)
		w       uint16
	)

	for i < n {
		w = endines.Uint16(p[i:])
		i += 2

		switch {
		case w <= 0x7f:
			out[j] = byte(w)
			j++
		case w <= 0x7ff:
			out[j] = byte(w>>6 | 0xc0)
			j++
			out[j] = byte(w&0x3f | 0x80)
			j++
		case w <= 0xffff:
			out[j] = byte(w>>12 | 0xe0)
			j++
			out[j] = byte(w>>6&0x3f | 0x80)
			j++
			out[j] = byte(w&0x3f | 0x80)
			j++
		}
	}
	_, err = dst.Write(out[:j])
	return err
}

// EncodeSlice encodes unicode to utf-8
func EncodeSlice(p []byte, endienes binary.ByteOrder) []byte {
	src := new(bytes.Buffer)
	src.Write(p)
	dst := new(bytes.Buffer)
	Encode(dst, src, endienes)
	return dst.Bytes()
}
