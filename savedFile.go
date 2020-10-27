package main

import (
	"io"
	"time"
)

type savedFile struct {
	Size uint64
	CTime, MTime time.Time
	CRC32              uint32
}

func (p *savedFile) readFileInfo  (r io.Reader) {
	p.Size = getUint64(r)
	p.CTime, p.MTime = time.Unix(getInt64(r),0), time.Unix( getInt64(r), 0)
	p.CRC32 = getDword(r)
}

func (p *savedFile) saveFileInfo  (w io.Writer) {
	saveUint64(w, p.Size)
	saveInt64(w, p.CTime.Unix())
	saveInt64(w, p.MTime.Unix())
	saveDword(w, p.CRC32)
}


