package main

import "io"

type dictSettings struct {
	SMorphology, SMorphFields string // rename!
	SStopwords                string // rename!
	Nwordforms                uint32
	IMinStemmingLen           uint32 // rename!
	BWordDict                 bool // rename!
	StopwordsUnstemmed        bool
	MorphFingerprint          string

	// embedded Stopwords
	EmbeddedStopwords bool
	Stopwords         []uint64
	StopwordFiles     []savedFile
	Stopwordsfnames   []string

	// embedded Wordforms
	EmbeddedWordforms bool
	Wordforms         []string
	WordformFiles     []savedFile
	Wordformfnames    []string
}

func (c *dictSettings) load(r io.Reader)  {
	c.SMorphology = getString(r)
	c.SMorphFields = getString(r)

	c.EmbeddedStopwords = getByteBool(r)
	if c.EmbeddedStopwords {
		nStopwords := getInt(r)
		c.Stopwords = make ([]uint64, nStopwords)
		for i:=0; i<nStopwords; i++ {
			c.Stopwords[i] = unzipOffsetBE(r)
		}
	}

	c.SStopwords = getString(r)

	nFiles:=getInt(r)
	c.StopwordFiles = make ([]savedFile,nFiles)
	c.Stopwordsfnames = make ([]string,nFiles)
	for i:=0; i<nFiles; i++ {
		c.Stopwordsfnames[i] = getString(r)
		c.StopwordFiles[i].readFileInfo(r)
	}

	c.EmbeddedWordforms = getByteBool(r)
	if c.EmbeddedWordforms {
		nWordforms := getInt(r)
		c.Wordforms = make ([]string,nWordforms)
		for i:=0; i<nWordforms; i++ {
			c.Wordforms[i] = getString(r)
		}
	}

	c.Nwordforms = getDword(r)

	c.WordformFiles = make ([]savedFile, c.Nwordforms)
	c.Wordformfnames = make ([]string, c.Nwordforms)
	for i:=0; i<int(c.Nwordforms); i++ {
		c.Wordformfnames[i] = getString(r)
		c.WordformFiles[i].readFileInfo(r)
	}

	c.IMinStemmingLen = getDword(r)
	c.BWordDict = getByteBool(r)
	c.StopwordsUnstemmed = getByteBool(r)
	c.MorphFingerprint = getString(r)
}

func (c *dictSettings) save(w io.Writer)  {
	saveString(w, c.SMorphology)
	saveString(w, c.SMorphFields)
	saveBoolByte(w, c.EmbeddedStopwords)
	if c.EmbeddedStopwords {
		nStopwords := len(c.Stopwords)
		saveInt(w, nStopwords)
		for i:=0; i<nStopwords; i++ {
			zipOffsetBE(w, c.Stopwords[i])
		}
	}
	saveString(w, c.SStopwords)

	nFiles:=len(c.StopwordFiles)
	saveInt(w,nFiles)
	for i:=0; i<nFiles; i++ {
		saveString(w, c.Stopwordsfnames[i])
		c.StopwordFiles[i].saveFileInfo(w)
	}

	saveBoolByte(w, c.EmbeddedWordforms)
	if c.EmbeddedWordforms {
		nWordforms := len(c.Wordforms)
		saveInt(w, nWordforms)
		for i:=0; i<nWordforms; i++ {
			saveString(w,c.Wordforms[i])
		}
	}

	saveDword(w, c.Nwordforms)
	for i:=0; i<int(c.Nwordforms); i++ {
		saveString(w, c.Wordformfnames[i])
		c.WordformFiles[i].saveFileInfo(w)
	}
	saveDword(w, c.IMinStemmingLen)
	saveBoolByte(w, c.BWordDict)
	saveBoolByte(w, c.StopwordsUnstemmed)
	saveString(w, c.MorphFingerprint)
}