package main

import "io"

type tokenizerSettings struct {
	IType                               byte
	SCaseFolding                        string
	IMinWordLen                         uint32
	SynonymsFile, Boundary, IgnoreChars string
	NgramLen                            uint32
	NgramChars, BlendChars, BlendMode   string

	// Synonyms
	EmbeddedSynonyms bool
	Synonyms         []string
	Synonymsfile     savedFile
}

func (c *tokenizerSettings) load(r io.Reader)  {
	c.IType = getByte(r)
	c.SCaseFolding = getString(r)
	c.IMinWordLen = getDword(r)

	c.EmbeddedSynonyms = getByteBool(r)
	if c.EmbeddedSynonyms {
		nSynonyms := getInt(r)
		c.Synonyms = make ([]string, nSynonyms)
		for i:=0; i<nSynonyms; i++ {
			c.Synonyms[i] = getString(r)
		}
	}
	c.SynonymsFile = getString(r)
	c.Synonymsfile.readFileInfo(r)
	c.Boundary, c.IgnoreChars = getString(r), getString(r)
	c.NgramLen = getDword(r)
	c.NgramChars, c.BlendChars, c.BlendMode = getString(r), getString(r), getString(r)
}

func (c *tokenizerSettings) save(w io.Writer)  {
	saveByte (w, c.IType)
	saveString(w, c.SCaseFolding)
	saveDword(w, c.IMinWordLen)

	saveBoolByte(w, c.EmbeddedSynonyms)
	if c.EmbeddedSynonyms {
		nSynonyms := len(c.Synonyms)
		saveInt(w, nSynonyms)
		for i:=0; i<nSynonyms; i++ {
			saveString(w, c.Synonyms[i])
		}
	}
	saveString(w, c.SynonymsFile)
	c.Synonymsfile.saveFileInfo(w)
	saveString(w, c.Boundary)
	saveString(w, c.IgnoreChars)
	saveDword(w, c.NgramLen)
	saveString(w, c.NgramChars)
	saveString(w, c.BlendChars)
	saveString(w, c.BlendMode)
}