package main

import "io"

type fieldFilterSettings struct {
	Regexps []string
}

func (c *fieldFilterSettings) load(r io.Reader)  {
	nfields := getInt(r)
	c.Regexps = make ([]string, nfields)
	if nfields==0 {
		return
	}
	for i:=0; i<nfields; i++ {
		c.Regexps[i] = getString(r)
	}
}

func (c *fieldFilterSettings) save(w io.Writer)  {
	nfields := len(c.Regexps)
	saveInt(w,nfields)
	if nfields==0 {
		return
	}
	for i:=0; i<nfields; i++ {
		saveString(w,c.Regexps[i])
	}
}
