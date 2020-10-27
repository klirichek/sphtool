package main

import (
	"fmt"
	"io"
	"os"
)

const MetaHeaderMagicPq uint32 = 0x50535451 ///< my magic 'PSTQ' header
const MetaVersionPq uint32 = 8

type metapq struct {
	metahdr
	IndexVersion        uint32
	Schema              indexschema
	IndexSettings       indexsettings
	TokSettings         tokenizerSettings
	DictSettings        dictSettings
	FieldFilterSettings fieldFilterSettings
	Queries             []storedQueryDesc
	Tid                 int64
}

func (m *metapq) load(r io.Reader) {
	// m.Hdr.load(r)
	m.IndexVersion = getDword(r)

	m.Schema.load(r, m.IndexVersion)
	m.IndexSettings.load(r, m.IndexVersion)
	m.TokSettings.load(r)
	m.DictSettings.load(r)
	if m.Version >= 6 {
		m.FieldFilterSettings.load(r)
	}

	nQueries := getDword(r)
	m.Queries = make([]storedQueryDesc, nQueries)
	for i := 0; i < int(nQueries); i++ {
		m.Queries[i].load(r, m.Version)
	}
	if m.Version > 6 {
		m.Tid = getInt64(r)
	}
}

func (m *metapq) save(w io.Writer) {
	m.metahdr.save(w)
	saveDword(w, m.IndexVersion)

	m.Schema.save(w, m.IndexVersion)
	m.IndexSettings.save(w, m.IndexVersion)
	m.TokSettings.save(w)
	m.DictSettings.save(w)
	if m.Version >= 6 {
		m.FieldFilterSettings.save(w)
	}
	nQueries := len(m.Queries)
	saveDword(w, uint32(nQueries))
	if nQueries == 0 {
		return
	}
	for i := 0; i < nQueries; i++ {
		m.Queries[i].save(w, m.Version)
	}
	if m.Version > 6 {
		saveUint64(w, i2u64(m.Tid))
	}
}

func (h *metahdr) isPq() bool {
	return h.Magic == MetaHeaderMagicPq
}

func (m *metapq) Checkvalid() {
	if m.Magic != MetaHeaderMagicPq {
		fmt.Printf("Wrong magic of the index: expected %d, got %d\n", MetaHeaderMagicPq, m.Magic)
		os.Exit(1)
	}
}
