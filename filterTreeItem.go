package main

import "io"

type filterTreeItem struct {
	Left  int32
	Right int32
	Item  int32
	Or    bool
}

func (h *filterTreeItem) loadV6(r io.Reader) {
	h.Left = getInt32(r)
	h.Right = getInt32(r)
	h.Item = getInt32(r)
	h.Or = getDword(r) != 0
}

func (h *filterTreeItem) saveV6(w io.Writer) {
	saveInt32(w, h.Left)
	saveInt32(w, h.Right)
	saveInt32(w, h.Item)
	uOr := uint32(0)
	if h.Or {
		uOr = 1
	}
	saveDword(w, uOr)
}

func (h *filterTreeItem) load(r io.Reader) {
	h.Left = int32(unzipDwordBE(r))
	h.Right = int32(unzipDwordBE(r))
	h.Item = int32(unzipDwordBE(r))
	h.Or = getByteBool(r)
}

func (h *filterTreeItem) save(w io.Writer) {
	zipDwordBE(w, i2u32(h.Left))
	zipDwordBE(w, i2u32(h.Right))
	zipDwordBE(w, i2u32(h.Item))
	uOr := uint32(0)
	if h.Or {
		uOr = 1
	}
	zipDwordBE(w, uOr)
}
