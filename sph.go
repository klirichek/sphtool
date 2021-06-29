package main

import (
	"fmt"
	"io"
	"os"
)

const IndexMagicHeader uint32 = 0x58485053
const SPHVersion uint32 = 63

type sph struct {
	metahdr
	Schema indexschema

	// dictionary header (wordlist checkpoints, infix blocks, etc)
	DictCheckpointsOffset uint64
	DictCheckPoints       uint32
	InfixCodepointBytes   byte
	InfixBlockOffset,
	InfixBlockWordSize uint32

	// index stats
	TotalDocuments                     uint32
	TotalBytes                         uint64
	IndexSettings                      indexsettings
	TokenizerSettings                  tokenizerSettings
	DictSettings                       dictSettings
	Docinfo, DocinfoIndex, MinMaxIndex uint64
	FieldLens                          []uint64
	FieldFilterSettings                fieldFilterSettings
}

func (h *sph) load(r io.Reader) {
	h.Magic = getDword(r)
	h.Version = getDword(r)
	h.Schema.load(r, h.Version)
	h.DictCheckpointsOffset, h.DictCheckPoints = getUint64(r), getDword(r)
	h.InfixCodepointBytes = getByte(r)
	h.InfixBlockOffset, h.InfixBlockWordSize = getDword(r), getDword(r)
	h.TotalDocuments = getDword(r)
	h.TotalBytes = getUint64(r)
	h.IndexSettings.load(r, h.Version)
	h.TokenizerSettings.load(r)
	h.DictSettings.load(r)
	h.Docinfo = getUint64(r)
	h.DocinfoIndex = getUint64(r)
	h.MinMaxIndex = getUint64(r)
	h.FieldFilterSettings.load(r)

	if h.IndexSettings.IndexFieldLens {
		lng := len(h.Schema.Fields)
		h.FieldLens = make([]uint64, lng)
		for i := 0; i < lng; i++ {
			h.FieldLens[i] = getUint64(r)
		}
	}
}

func (h *sph) save(w io.Writer) {

	saveDword(w, h.Magic)
	saveDword(w, h.Version)
	h.Schema.save(w, h.Version)
	saveUint64(w, h.DictCheckpointsOffset)
	saveDword(w, h.DictCheckPoints)
	saveByte(w, h.InfixCodepointBytes)
	saveDword(w, h.InfixBlockOffset)
	saveDword(w, h.InfixBlockWordSize)
	saveDword(w, h.TotalDocuments)
	saveUint64(w, h.TotalBytes)
	h.IndexSettings.save(w, h.Version)
	h.TokenizerSettings.save(w)
	h.DictSettings.save(w)
	saveUint64(w, h.Docinfo)
	saveUint64(w, h.DocinfoIndex)
	saveUint64(w, h.MinMaxIndex)
	h.FieldFilterSettings.save(w)

	if h.IndexSettings.IndexFieldLens {
		lng := len(h.Schema.Fields)
		h.FieldLens = make([]uint64, lng)
		for i := 0; i < lng; i++ {
			saveUint64(w, h.FieldLens[i])
		}
	}
}

func (h *metahdr) isSph() bool {
	return h.Magic == IndexMagicHeader
}

func (h *sph) Checkvalid() {
	if h.Magic != IndexMagicHeader {
		fmt.Printf("Wrong magic of the index: expected %d, got %d\n", IndexMagicHeader, h.Magic)
		os.Exit(1)
	}
}
