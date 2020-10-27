package main

import "io"

type storedQueryDesc struct {
	QUID       int64
	Query      string
	QL         bool
	Tags       string
	Filters    []filterSettings
	FilterTree []filterTreeItem
}

// ver 1
func (c *storedQueryDesc) loadV1(r io.Reader) {
	c.Query = getString(r)
}

func (c *storedQueryDesc) saveV1(w io.Writer) {
	saveString(w, c.Query)
}

// ver 2
func (c *storedQueryDesc) loadV2(r io.Reader) {
	c.Query = getString(r)
	c.Tags = getString(r)

	nFilters := getDword(r)
	nTree := getDword(r)
	c.Filters = make([]filterSettings, nFilters)
	c.FilterTree = make([]filterTreeItem, nTree)
	for i := 0; i < int(nFilters); i++ {
		c.Filters[i].loadV6(r)
	}
	for i := 0; i < int(nTree); i++ {
		c.FilterTree[i].loadV6(r)
	}
}

func (c *storedQueryDesc) saveV2(w io.Writer) {
	saveString(w, c.Query)
	saveString(w, c.Tags)

	nFilters := len(c.Filters)
	nTree := len(c.FilterTree)
	saveInt(w, nFilters)
	saveInt(w, nTree)
	for i := 0; i < nFilters; i++ {
		c.Filters[i].saveV6(w)
	}
	for i := 0; i < nTree; i++ {
		c.FilterTree[i].saveV6(w)
	}
}

// ver 3
func (c *storedQueryDesc) loadV3(r io.Reader) {
	c.QUID = getInt64(r)
	c.Query = getString(r)
	c.Tags = getString(r)

	nFilters := getDword(r)
	nTree := getDword(r)
	c.Filters = make([]filterSettings, nFilters)
	c.FilterTree = make([]filterTreeItem, nTree)
	for i := 0; i < int(nFilters); i++ {
		c.Filters[i].loadV6(r)
	}
	for i := 0; i < int(nTree); i++ {
		c.FilterTree[i].loadV6(r)
	}
}

func (c *storedQueryDesc) saveV3(w io.Writer) {
	saveInt64(w, c.QUID)
	saveString(w, c.Query)
	saveString(w, c.Tags)

	nFilters := len(c.Filters)
	nTree := len(c.FilterTree)
	saveInt(w, nFilters)
	saveInt(w, nTree)
	for i := 0; i < nFilters; i++ {
		c.Filters[i].saveV6(w)
	}
	for i := 0; i < nTree; i++ {
		c.FilterTree[i].saveV6(w)
	}
}

// ver 4
func (c *storedQueryDesc) loadV4(r io.Reader) {
	c.QUID = getInt64(r)
	c.QL = getDword(r) != 0
	c.Query = getString(r)
	c.Tags = getString(r)

	nFilters := getDword(r)
	nTree := getDword(r)
	c.Filters = make([]filterSettings, nFilters)
	c.FilterTree = make([]filterTreeItem, nTree)
	for i := 0; i < int(nFilters); i++ {
		c.Filters[i].loadV6(r)
	}
	for i := 0; i < int(nTree); i++ {
		c.FilterTree[i].loadV6(r)
	}
}

func (c *storedQueryDesc) saveV4(w io.Writer) {
	saveInt64(w, c.QUID)
	saveBoolDword(w, c.QL)
	saveString(w, c.Query)
	saveString(w, c.Tags)

	nFilters := len(c.Filters)
	nTree := len(c.FilterTree)
	saveInt(w, nFilters)
	saveInt(w, nTree)
	for i := 0; i < nFilters; i++ {
		c.Filters[i].saveV6(w)
	}
	for i := 0; i < nTree; i++ {
		c.FilterTree[i].saveV6(w)
	}
}

// ver 6
func (c *storedQueryDesc) loadDefault(r io.Reader) {
	c.QUID = u2i64(unzipOffsetBE(r))
	c.QL = getByteBool(r)
	c.Query = getString(r)
	c.Tags = getString(r)

	nFilters := unzipDwordBE(r)
	nTree := unzipDwordBE(r)
	var i uint32
	c.Filters = make([]filterSettings, nFilters)
	c.FilterTree = make([]filterTreeItem, nTree)
	for i = 0; i < nFilters; i++ {
		c.Filters[i].load(r)
	}
	for i = 0; i < nTree; i++ {
		c.FilterTree[i].load(r)
	}
}

func (c *storedQueryDesc) saveDefault(w io.Writer) {
	zipOffsetBE(w, i2u64(c.QUID))
	saveBoolByte(w, c.QL)
	saveString(w, c.Query)
	saveString(w, c.Tags)
	var nFilters, nTree, i uint32
	nFilters = uint32(len(c.Filters))
	nTree = uint32(len(c.FilterTree))
	zipDwordBE(w, nFilters)
	zipDwordBE(w, nTree)
	for i = 0; i < nFilters; i++ {
		c.Filters[i].save(w)
	}
	for i = 0; i < nTree; i++ {
		c.FilterTree[i].save(w)
	}
}

func (c *storedQueryDesc) load(r io.Reader, uVersion uint32) {
	switch uVersion {
	case 1:
		c.loadV1(r)
	case 2:
		c.loadV2(r)
	case 3:
		c.loadV3(r)
	case 4:
		c.loadV4(r)
	}
	c.loadDefault(r)
}

func (c *storedQueryDesc) save(w io.Writer, uVersion uint32) {
	switch uVersion {
	case 1:
		c.saveV1(w)
	case 2:
		c.saveV2(w)
	case 3:
		c.saveV3(w)
	case 4:
		c.saveV4(w)
	}
	c.saveDefault(w)
}
