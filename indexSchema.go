package main

import "io"

type columnInfo struct {
	Name                string
	Etype               uint32
	CompatRowItem		uint32
	Bitoffset, Bitcount int32
	Payload             bool
}

func readCol(rd io.Reader, uVersion uint32) columnInfo {
	var c columnInfo
	c.Name = getString(rd)
	c.Etype = getDword(rd)
	if uVersion<57 {
		c.CompatRowItem = getDword(rd)
		c.Bitoffset = getInt32(rd)
		c.Bitcount = getInt32(rd)
	}
	c.Payload = getByteBool(rd)
	return c
}

func (c *columnInfo )saveCol (w io.Writer, uVersion uint32) {
	saveString(w, c.Name)
	saveDword(w, c.Etype)
	if uVersion<57 {
		saveDword(w,c.CompatRowItem)
		saveInt(w, int(c.Bitoffset))
		saveInt(w, int(c.Bitcount))
	}
	saveBoolByte(w, c.Payload)
}

type indexschema struct {
	Fields, Attrs []columnInfo
}

func (c *indexschema) load(rd io.Reader, uVersion uint32) {
	nfields := int(getDword(rd))
	for i:=0; i<nfields; i++ {
		c.Fields = append (c.Fields, readCol(rd,uVersion))
	}

	nattrs := int(getDword(rd))
	for i:=0; i<nattrs; i++ {
		c.Attrs = append (c.Attrs, readCol(rd,0))
	}
}

func (c *indexschema) save(w io.Writer, uVersion uint32) {
	nfields := len(c.Fields)
	saveInt(w,nfields)
	for i:=0; i<nfields; i++ {
		c.Fields[i].saveCol(w,uVersion)
	}

	nattrs := len(c.Attrs)
	saveInt(w,nattrs)
	for i:=0; i<nattrs; i++ {
		c.Attrs[i].saveCol(w,0)
	}
}
