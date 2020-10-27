package main

import (
	"fmt"
	"io"
	"os"
)

const MetaHeaderMagic uint32 = 0x54525053 ///< my magic 'SPRT' header
const MetaVersion uint32 = 17

type metahdr struct {
	Magic   uint32
	Version uint32
}

type meta struct {
	metahdr
	TotalDocuments uint32
	TotalBytes     uint64
	TID            uint64
	SettingsVer    uint32

	Schema              indexschema
	IndexSettings       indexsettings
	TokSettings         tokenizerSettings
	DictSettings        dictSettings
	WordsCheckpoints    uint32
	MaxCodepointLength  uint32
	BloomKeyLen         byte
	BloomHashesCount    byte
	FieldFilterSettings fieldFilterSettings
	ChunkNames          []int32
	SoftRamLimit        int64
}

func (h *meta) loadChunkNames(r io.Reader) {
	nfields := getInt(r)
	h.ChunkNames = make([]int32, nfields)
	if nfields == 0 {
		return
	}
	for i := 0; i < nfields; i++ {
		h.ChunkNames[i] = getInt32(r)
	}
}

func (h *meta) saveChunkNames(w io.Writer) {
	nfields := len(h.ChunkNames)
	saveInt(w, nfields)
	if nfields == 0 {
		return
	}
	for i := 0; i < nfields; i++ {
		saveInt32(w, h.ChunkNames[i])
	}
}

func (h *metahdr) load(r io.Reader) {
	h.Magic = getDword(r)
	h.Version = getDword(r)
}

func (h *metahdr) save(w io.Writer) {
	saveDword(w, h.Magic)
	saveDword(w, h.Version)
}

func (h *meta) load(r io.Reader) {
	// h.Hdr.load(r)
	h.TotalDocuments = getDword(r)
	h.TotalBytes = getUint64(r)
	h.TID = getUint64(r)
	h.SettingsVer = getDword(r)

	h.Schema.load(r, h.SettingsVer)
	h.IndexSettings.load(r, h.SettingsVer)
	h.TokSettings.load(r)
	h.DictSettings.load(r)

	h.WordsCheckpoints = getDword(r)
	h.MaxCodepointLength = getDword(r)
	h.BloomKeyLen = getByte(r)
	h.BloomHashesCount = getByte(r)

	h.FieldFilterSettings.load(r)
	h.loadChunkNames(r)

	if h.Version >= 17 {
		h.SoftRamLimit = getInt64(r)
	}
}

func (h *meta) save(w io.Writer) {
	h.metahdr.save(w)
	saveDword(w, h.TotalDocuments)
	saveUint64(w, h.TotalBytes)
	saveUint64(w, h.TID)
	saveDword(w, h.SettingsVer)
	h.Schema.save(w, h.SettingsVer)
	h.IndexSettings.save(w, h.SettingsVer)
	h.TokSettings.save(w)
	h.DictSettings.save(w)

	saveDword(w, h.WordsCheckpoints)
	saveDword(w, h.MaxCodepointLength)
	saveByte(w, h.BloomKeyLen)
	saveByte(w, h.BloomHashesCount)

	h.FieldFilterSettings.save(w)
	h.saveChunkNames(w)

	if h.Version >= 17 {
		saveInt64(w, h.SoftRamLimit)
	}
}

func (h *metahdr) isRt() bool {
	return h.Magic == MetaHeaderMagic
}

func (h *meta) Checkvalid() {
	if h.Magic != MetaHeaderMagic {
		fmt.Printf("Wrong magic of the index: expected %d, got %d\n", MetaHeaderMagic, h.Magic)
		os.Exit(1)
	}
}
