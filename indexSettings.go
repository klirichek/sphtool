package main

import "io"

type indexsettings struct {
	MinPrefixLen, MinInfixLen, MaxSubstringLen               int
	HtmlStrip                                                bool
	HtmlIndexAttrs, HtmlRemoveElements                       string
	IndexExactWords                                          bool
	EHitless, EHitFormat                                     uint32
	IndexSp                                                  bool
	Zones                                                    string
	BoundaryStep, StopwordStep, OvershortStep, EmbeddedLimit int
	EBigramIndex                                             byte
	BigramWords                                              string
	IndexFieldLens                                           bool
	EChineseRLP                                              byte
	RLPContext                                               string
	IndexTokenFilter                                         string
	BlobUpdateSpace                                          uint64
	SkiplistBlockSize                                        uint32
	HitlessFiles											 string
}

func (c *indexsettings) load (r io.Reader, ver uint32) {
	c.MinPrefixLen, c.MinInfixLen, c.MaxSubstringLen = getInt(r), getInt(r), getInt(r)
	c.HtmlStrip = getByteBool(r)
	c.HtmlIndexAttrs, c.HtmlRemoveElements = getString(r), getString(r)
	c.IndexExactWords = getByteBool(r)
	c.EHitless, c.EHitFormat = getDword(r), getDword(r)
	c.IndexSp = getByteBool(r)
	c.Zones = getString(r)
	c.BoundaryStep, c.StopwordStep, c.OvershortStep, c.EmbeddedLimit = getInt(r), getInt(r), getInt(r), getInt(r)
	c.EBigramIndex = getByte(r)
	c.BigramWords = getString(r)
	c.IndexFieldLens = getByteBool(r)
	c.EChineseRLP = getByte(r)
	c.RLPContext = getString(r)
	c.IndexTokenFilter = getString(r)
	c.BlobUpdateSpace = getUint64(r)
	c.SkiplistBlockSize = 128
	if ver>=56 {
		c.SkiplistBlockSize = getDword(r)
	}
	if ver>=60 {
		c.HitlessFiles = getString(r)
	}
}

func (c *indexsettings) save(w io.Writer, ver uint32) {
	saveInt(w, c.MinPrefixLen)
	saveInt(w, c.MinInfixLen)
	saveInt(w, c.MaxSubstringLen)
	saveBoolByte(w, c.HtmlStrip)
	saveString(w, c.HtmlIndexAttrs)
	saveString(w, c.HtmlRemoveElements)
	saveBoolByte(w, c.IndexExactWords)
	saveDword(w, c.EHitless)
	saveDword(w, c.EHitFormat)
	saveBoolByte(w, c.IndexSp)
	saveString(w, c.Zones)
	saveInt(w, c.BoundaryStep)
	saveInt(w, c.StopwordStep)
	saveInt(w, c.OvershortStep)
	saveInt(w, c.EmbeddedLimit)
	saveByte(w, c.EBigramIndex)
	saveString(w, c.BigramWords)
	saveBoolByte(w, c.IndexFieldLens)
	saveByte(w, c.EChineseRLP)
	saveString(w, c.RLPContext)
	saveString(w, c.IndexTokenFilter)
	saveUint64(w, c.BlobUpdateSpace)
	if ver>=56 {
		saveDword(w, c.SkiplistBlockSize)
	}
	if ver>=60 {
		saveString(w, c.HitlessFiles)
	}
}
