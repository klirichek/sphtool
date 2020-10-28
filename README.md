# SPHtool - render/parse sph/meta files

Manticore search indexes consists from huge datafiles (as dictionaries, attributes, etc) and small headers containing
basic meta info as index settings, paths, etc. That is binary internal format and usually not capable to be edited from
outside (so consider it as 'black box').

This tool renders whole header to json (which is obviously much more human-friendly).
Also it performs opposite operation - parses json back to header. 
 
    ./sphtool -i index.sph
    ./sphtool -i index.meta
    
parses `index.sph` and produce `index.sph.json` as output. The same way works with 'meta' headers (both from realtime
and percolate indexes). File extension (`.sph` or `.meta`) right now used to determine file format and is mandatory
(I know, file fingerprint is enouth go determine format, however right now extension is used). By default target file is
original one plus `.json` extension.

    ./sphtool -i index.sph.json
    
renders `index.sph.json` and produce `index.sph` as output. File extension `.json` is mandatory. Internal flavour (sph,
meta from rt or saved percolate queries) is determined from magic field in json. By default target file is original one
without `.json` extension.

Optionally target file may be provided with `-o` option, as:

    ./sphtool -i pqindex.json -o pqindex.meta
    ./sphtool -i plainindex.json -o plainindex.sph

Supported versions of headers are visible in banner:

    SPHtool v $Format:%h$ - sph and meta files converter
    Copyright (c) 2020 Alexey N. Vinogradov (a.n.vinogradov@gmail.com)
    Supports sph v. 60, meta v. 17, pq meta v. 8
    Usage of ./sphtool:
      -i string
        	path to sph, meta, sph.json or meta.json file
      -o string
        	path to target file (optional)
        	
Tool may convert to json and back following headers:
 * `.sph` - which is 'sphinx header' - main file containing plain index params as schema, tokenenizer settings, etc.
 That is both part of standalone plain  indexes and rt disk chunks. In json you may change, say, path to files 
 (wordforms, stopwords, whatever) if they changed.  Or fixup timestamps of the files, as they used to check files 
 consistency.
 * `.meta` of rt - is header of realtime index. It has similar sense as sph file of plain index, but in aggregate
 meaning over all index chunks.
 * `.meta` of pq - is actually whole pq index (it has nothing else apart this meta). Apart general settings (like schema)
 and embedded files/hashes it includes stored queries.
 
 ## Building
 
 Use
 
    go get github.com/klirichek/sphtool
    
or (if you have clone sources):

    cd /src/of/sphtool
    go build -o sphtool . 

## Caveats

It might look useful to format json if you want to edit it.
Many sources advice to use 'jq' tool for it, like:

    cat index.meta.json | jq . > index_formatted.meta.json
    
However that `jq` corrupts some numbers when reformat, so that final file became unusable. Look at this live
example (one and same file first printed, then passed to jq). Note the difference! BEWARE! 

    $ cat ex.json
    {
      "numbers": [
        623237308590442816,
        2696721519739033356,
        15902905282948881040
      ]
    }
    
    $ cat ex.json | jq .
    {
      "numbers": [
        623237308590442800,
        2696721519739033600,
        15902905282948880000
      ]
    }
    

 