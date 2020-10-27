package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"
)

func checkl(i int, e error) {
	if e != nil {
		fmt.Printf("Couldn't read %d bytes\n", i)
		//		panic(e)
	}
}

func getByte(rd io.Reader) byte {
	bf := make([]byte, 1)
	lng, err := rd.Read(bf)
	checkl(lng, err)
	return bf[0]
}

func getDword(rd io.Reader) uint32 {
	bf := make([]byte, 4)
	lng, err := rd.Read(bf)
	checkl(lng, err)
	return binary.LittleEndian.Uint32(bf)
}

func getInt(rd io.Reader) int {
	return int(getDword(rd))
}

func getInt32(rd io.Reader) int32 {
	return int32(getDword(rd))
}

func getUint64(rd io.Reader) uint64 {
	bf := make([]byte, 8)
	lng, err := rd.Read(bf)
	checkl(lng, err)
	return binary.LittleEndian.Uint64(bf)
}

func getInt64(rd io.Reader) int64 {
	return u2i64(getUint64(rd))
}

// VBL BE decoding (most significant bits encoded first)
func unzipOffsetBE(rd io.Reader) uint64 {
	var res uint64 = 0
	for true {
		b := getByte(rd)
		res = res<<7 + uint64(b&0x7f)
		if b&0x80 == 0 {
			break
		}
	}
	return res
}

func u2i64(u uint64) int64 {
	return *(*int64)(unsafe.Pointer(&u))
}

func i2u64(u int64) uint64 {
	return *(*uint64)(unsafe.Pointer(&u))
}

func unzipDwordBE(rd io.Reader) uint32 {
	var res uint32 = 0
	for true {
		b := getByte(rd)
		res = res<<7 + uint32(b&0x7f)
		if b&0x80 == 0 {
			break
		}
	}
	return res
}

func i2u32(u int32) uint32 {
	return *(*uint32)(unsafe.Pointer(&u))
}

// VBL LE decoding (most significant bits encoded last)
func getByteBool(rd io.Reader) bool {
	return getByte(rd) != 0
}

func getString(rd io.Reader) string {
	lng := getInt(rd)
	result := make([]byte, lng)
	if lng > 0 {
		lng, err := rd.Read(result)
		checkl(lng, err)
	}
	return string(result)
}

func saveByte(w io.Writer, v byte) {
	bf := make([]byte, 1)
	bf[0] = v
	lng, err := w.Write(bf)
	checkl(lng, err)
}

func saveDword(w io.Writer, v uint32) {
	bf := make([]byte, 4)
	binary.LittleEndian.PutUint32(bf, v)
	lng, err := w.Write(bf)
	checkl(lng, err)
}

func saveInt(w io.Writer, v int) {
	saveDword(w, uint32(v))
}

func saveInt32(w io.Writer, v int32) {
	saveDword(w, uint32(v))
}

func saveUint64(w io.Writer, v uint64) {
	bf := make([]byte, 8)
	binary.LittleEndian.PutUint64(bf, v)
	lng, err := w.Write(bf)
	checkl(lng, err)
}

func saveInt64(w io.Writer, v int64) {
	saveUint64(w, i2u64(v))
}

func calcZippedLen(v uint64) int {
	n := 0
	for true {
		n++
		v >>= 7
		if v == 0 {
			break
		}
	}
	return n
}

// VBL BE encoding (most significant bits encoded first)
func zipOffsetBE(w io.Writer, v uint64) {
	bytes := calcZippedLen(v)
	o := make([]byte, bytes)
	for i := 0; i < bytes; i++ {
		o[i] = byte(v >> (7 * (bytes - i - 1)) & 0x7F)
	}
	for i := 0; i < len(o)-1; i++ {
		o[i] |= 0x80
	}
	lng, err := w.Write(o)
	checkl(lng, err)
}

func zipDwordBE(w io.Writer, v uint32) {
	zipOffsetBE(w, uint64(v))
}

func saveBoolByte(w io.Writer, v bool) {
	if v {
		saveByte(w, 1)
	} else {
		saveByte(w, 0)
	}
}

func saveBoolDword(w io.Writer, v bool) {
	if v {
		saveDword(w, 1)
	} else {
		saveDword(w, 0)
	}
}

func saveString(w io.Writer, v string) {
	saveInt(w, len(v))
	if len(v) > 0 {
		lng, err := w.Write([]byte(v))
		checkl(lng, err)
	}
}
