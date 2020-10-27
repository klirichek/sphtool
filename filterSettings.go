package main

import (
	"io"
)

type filterSettings struct {
	AttrName    string
	Exclude     bool
	HasEqualMin bool
	HasEqualMax bool
	EType       uint32
	MvaFunc     uint32
	MinValue    int64
	MaxValue    int64
	Values      []int64
	Strings     []string
	OpenLeft    bool
	OpenRight   bool
	IsNull      bool
}

func (c *filterSettings) loadV6(r io.Reader) {
	c.AttrName = getString(r)
	c.Exclude = getDword(r) != 0
	c.HasEqualMin = getDword(r) != 0
	c.HasEqualMax = getDword(r) != 0
	c.EType = getDword(r)
	c.MvaFunc = getDword(r)

	c.MinValue = getInt64(r)
	c.MaxValue = getInt64(r)

	nValues := getDword(r)
	nStrings := getDword(r)
	c.Values = make([]int64, nValues)
	c.Strings = make([]string, nStrings)
	for i := 0; i < int(nValues); i++ {
		c.Values[i] = getInt64(r)
	}
	for i := 0; i < int(nStrings); i++ {
		c.Strings[i] = getString(r)
	}
}

func (c *filterSettings) saveV6(w io.Writer) {
	saveString(w, c.AttrName)
	saveBoolDword(w, c.Exclude)
	saveBoolDword(w, c.HasEqualMin)
	saveBoolDword(w, c.HasEqualMax)
	saveDword(w, c.EType)
	saveDword(w, c.MvaFunc)

	saveInt64(w, c.MinValue)
	saveInt64(w, c.MaxValue)

	nValues := len(c.Values)
	nStrings := len(c.Strings)
	saveInt(w, nValues)
	saveInt(w, nStrings)
	for i := 0; i < nValues; i++ {
		saveInt64(w, c.Values[i])
	}
	for i := 0; i < nStrings; i++ {
		saveString(w, c.Strings[i])
	}
}

func (c *filterSettings) load(r io.Reader) {
	c.AttrName = getString(r)
	c.Exclude = unzipDwordBE(r) != 0
	c.HasEqualMin = unzipDwordBE(r) != 0
	c.HasEqualMax = unzipDwordBE(r) != 0
	c.OpenLeft = unzipDwordBE(r) != 0
	c.OpenRight = unzipDwordBE(r) != 0
	c.IsNull = unzipDwordBE(r) != 0
	c.EType = unzipDwordBE(r)
	c.MvaFunc = unzipDwordBE(r)
	c.MinValue = u2i64(unzipOffsetBE(r))
	c.MaxValue = u2i64(unzipOffsetBE(r))
	var nValues, nStrings, i uint32
	nValues = unzipDwordBE(r)
	nStrings = unzipDwordBE(r)
	c.Values = make([]int64, nValues)
	c.Strings = make([]string, nStrings)
	for i = 0; i < nValues; i++ {
		c.Values[i] = u2i64(unzipOffsetBE(r))
	}
	for i = 0; i < nStrings; i++ {
		c.Strings[i] = getString(r)
	}
}

func (c *filterSettings) save(w io.Writer) {
	saveString(w, c.AttrName)
	saveBoolByte(w, c.Exclude)
	saveBoolByte(w, c.HasEqualMin)
	saveBoolByte(w, c.HasEqualMax)
	saveBoolByte(w, c.OpenLeft)
	saveBoolByte(w, c.OpenRight)
	saveBoolByte(w, c.IsNull)
	zipDwordBE(w, c.EType)
	zipDwordBE(w, c.MvaFunc)
	zipOffsetBE(w, i2u64(c.MinValue))
	zipOffsetBE(w, i2u64(c.MaxValue))
	var nValues, nStrings, i uint32
	nValues = uint32(len(c.Values))
	nStrings = uint32(len(c.Strings))
	zipDwordBE(w, nValues)
	zipDwordBE(w, nStrings)
	for i = 0; i < nValues; i++ {
		zipOffsetBE(w, i2u64(c.Values[i]))
	}
	for i = 0; i < nStrings; i++ {
		saveString(w, c.Strings[i])
	}
}
