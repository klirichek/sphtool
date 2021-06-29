package main

import "io"

type AttrEngine_e uint32

const (
	DEFAULT = AttrEngine_e(iota)
	ROWWISE
	COLUMNAR
)

type AttrFlags_e uint32

const (
	ATTR_NONE            = 0
	ATTR_COLUMNAR        = 1
	ATTR_COLUMNAR_HASHES = 2
)

type FieldFlags_e uint32

const (
	FIELD_NONE         = 0
	FIELD_STORED       = 1
	FIELD_INDEXED      = 2
	FIELD_IS_ATTRIBUTE = 4
)

type columnInfo struct {
	Name                string
	Etype               uint32
	CompatRowItem       uint32
	Bitoffset, Bitcount int32
	FieldFlags          FieldFlags_e
	AttrFlags           AttrFlags_e
	AttrEngine          AttrEngine_e
	Payload             bool
}

func readField(rd io.Reader, uVersion uint32) columnInfo {
	var c columnInfo
	if uVersion >= 57 {
		c.Name = getString(rd)
		c.Etype = getDword(rd)
		c.Payload = getByteBool(rd)
	} else {
		c = readCol(rd, uVersion)
	}
	if uVersion < 59 {
		c.FieldFlags |= FIELD_INDEXED
	}
	return c
}

func (c *columnInfo) saveField(w io.Writer) {
	saveString(w, c.Name)
	saveDword(w, c.Etype)
	saveBoolByte(w, c.Payload)
}

func readCol(rd io.Reader, uVersion uint32) columnInfo {
	var c columnInfo
	c.Name = getString(rd)
	if c.Name == "" {
		c.Name = "@emptyname"
	}
	c.Etype = getDword(rd)
	c.CompatRowItem = getDword(rd)
	c.Bitoffset = getInt32(rd)
	c.Bitcount = getInt32(rd)
	c.Payload = getByteBool(rd)

	if uVersion >= 61 {
		c.AttrFlags = AttrFlags_e(getDword(rd))
	}

	if uVersion >= 63 {
		c.AttrEngine = AttrEngine_e(getDword(rd))
	}

	return c
}

func (c *columnInfo) saveCol(w io.Writer, uVersion uint32) {
	saveString(w, c.Name)
	saveDword(w, c.Etype)
	saveDword(w, c.CompatRowItem)
	saveInt(w, int(c.Bitoffset))
	saveInt(w, int(c.Bitcount))
	saveBoolByte(w, c.Payload)
	if uVersion >= 61 {
		saveDword(w, uint32(c.AttrFlags))
	}

	if uVersion >= 63 {
		saveDword(w, uint32(c.AttrEngine))
	}
}

type indexschema struct {
	Fields, Attrs []columnInfo
}

func (c *indexschema) load(rd io.Reader, uVersion uint32) {
	nfields := int(getDword(rd))
	for i := 0; i < nfields; i++ {
		c.Fields = append(c.Fields, readField(rd, uVersion))
	}

	nattrs := int(getDword(rd))
	for i := 0; i < nattrs; i++ {
		c.Attrs = append(c.Attrs, readCol(rd, uVersion))
	}
}

func (c *indexschema) save(w io.Writer, uVersion uint32) {
	nfields := len(c.Fields)
	saveInt(w, nfields)
	for i := 0; i < nfields; i++ {
		c.Fields[i].saveField(w)
	}

	nattrs := len(c.Attrs)
	saveInt(w, nattrs)
	for i := 0; i < nattrs; i++ {
		c.Attrs[i].saveCol(w, uVersion)
	}
}
