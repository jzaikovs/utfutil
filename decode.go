package utfutil

import (
	"encoding/binary"
	"io"
	"io/ioutil"
)

func Decode(dst io.Writer, src io.Reader) error {
	p, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}

	var endines binary.ByteOrder

	if p[0] == 0xef && p[1] == 0xbb && p[2] == 0xbf {
		endines = binary.LittleEndian
	} else {
		endines = binary.BigEndian
		if _, err = dst.Write([]byte{0xfe, 0xff}); err != nil {
			return err
		}
		p = p[3:]
	}

	i := 0
	j := 0
	out := make([]byte, len(p)*2)
	var w uint16
	for i < len(p) {
		switch {
		case p[i]&0xf0 == 0xe0:
			//3bytes
			w = uint16(p[i]) & 0xf << 12
			i++
			w |= uint16(p[i]) & 0x3f << 6
			i++
			w |= uint16(p[i]) & 0x3f
			i++
		case p[i]&0xe0 == 0xc0:
			//2byte
			w = uint16(p[i]) & 0x1f << 6
			i++
			w |= uint16(p[i]) & 0x3f
			i++
		case p[i]&0x80 == 0:
			//1byte
			w = uint16(p[i])
			i++
		}

		endines.PutUint16(out[j:], w)
		j += 2
	}
	_, err = dst.Write(out[:j])
	return err
}
