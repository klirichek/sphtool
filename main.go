package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

type actionFlags uint32

const (
	unknown actionFlags = iota
	sph2json
	meta2json
	json2
)

func getAction(inf string) (actionFlags, string, error) {
	ln := len(inf)
	if ln < 3 {
		return unknown, "", errors.New(fmt.Sprintf("wrong file provided (%s)", inf))
	}

	var outfile string
	if inf[len(inf)-3:] == "sph" {
		outfile = inf + ".json"
		return sph2json, outfile, nil
	}

	if inf[len(inf)-4:] == "meta" {
		outfile = inf + ".json"
		return meta2json, outfile, nil
	}

	if ln < 5 || inf[len(inf)-5:] != ".json" {
		return unknown, "", errors.New(fmt.Sprintf("wrong file provided (%s)", inf))
	}

	outfile = inf[0 : len(inf)-5]
	return json2, outfile, nil
}

func check(e error) {
	if e != nil {
		fmt.Printf("%v", e)
		os.Exit(1)
	}
}

func main() {

	inFile := flag.String("i", "", "path to sph, meta, sph.json or meta.json file")
	poutf := flag.String("o", "", "path to target file (optional)")

	flag.Parse()

	version := "$Format:%h$"

	if inFile == nil || *inFile == "" {
		fmt.Printf("SPHtool v %s - sph and meta files converter\n", version)
		fmt.Printf("Copyright (c) 2020-2021 Alexey N. Vinogradov (a.n.vinogradov@gmail.com)\n")
		fmt.Printf("Supports sph v. %d, meta v. %d, pq meta v. %d\n", SPHVersion, MetaVersion, MetaVersionPq)
		flag.Usage()
		os.Exit(1)
	}

	action, outf, err := getAction(*inFile)
	check(err)

	if poutf != nil && *poutf != "" {
		outf = *poutf
	}

	r, err := os.Open(*inFile)
	check(err)
	defer func() { _ = r.Close() }()

	o, err := os.Create(outf)
	check(err)
	defer func() { _ = o.Close() }()

	switch action {
	case sph2json:
		fmt.Printf("Input sph file is %v\n", *inFile)
		var header sph
		header.load(r)
		header.Checkvalid() // will abort inside if this fails
		enc := json.NewEncoder(o)
		err = enc.Encode(&header)
		check(err)

		fmt.Printf("File serialized to json %v\n", outf)

	case meta2json:
		fmt.Printf("Input meta file is %v\n", *inFile)
		enc := json.NewEncoder(o)
		var hdr metahdr
		var err error
		hdr.load(r)
		if hdr.isPq() {
			var pqmeta metapq
			pqmeta.metahdr = hdr
			pqmeta.load(r)
			pqmeta.Checkvalid() // will abort inside if this fails
			err = enc.Encode(&pqmeta)
		} else if hdr.isRt() {
			var rtmeta meta
			rtmeta.metahdr = hdr
			rtmeta.load(r)
			rtmeta.Checkvalid() // will abort inside if this fails
			err = enc.Encode(&rtmeta)
		}
		check(err)
		fmt.Printf("File serialized to json %v\n", outf)

	case json2:
		fmt.Printf("Input json file is %v\n", *inFile)
		var hdr metahdr
		{
			dec := json.NewDecoder(r)
			err = dec.Decode(&hdr)
			check(err)
		}
		_, _ = r.Seek(0, 0)
		if hdr.isSph() {
			var header sph
			dec := json.NewDecoder(r)
			err = dec.Decode(&header)
			check(err)
			header.save(o)
			fmt.Printf("Json encoded to sph %v\n", outf)
		} else if hdr.isRt() {
			var header meta
			dec := json.NewDecoder(r)
			err = dec.Decode(&header)
			check(err)
			header.save(o)
			fmt.Printf("Json encoded to meta %v\n", outf)
		} else if hdr.isPq() {
			var header metapq
			dec := json.NewDecoder(r)
			err = dec.Decode(&header)
			check(err)
			header.save(o)
			fmt.Printf("Json encoded to pq meta %v\n", outf)
		}
	}
}
