# logstash-config : parser and abstract syntax tree for [Logstash](https://github.com/elastic/logstash) config files

[![GoDoc](https://godoc.org/github.com/breml/logstash-config?status.svg)](https://godoc.org/github.com/breml/logstash-config) [![License](https://img.shields.io/badge/license-Apache_2.0-blue.svg)](LICENSE)

## Overview

The Go package config provides a ready to use parser for [Logstash](https://github.com/elastic/logstash) configuration files.

The basis of the grammar for the parsing of the Logstash configuration format is the original [Logstash Treetop grammar](https://github.com/elastic/logstash/blob/master/logstash-core/lib/logstash/config/grammar.treetop) which could be used with only minor changes.

logstash-config uses [pigeon](https://github.com/mna/pigeon) to generate the parser from the PEG (parser expression grammar). Special thanks to Martin Angers ([mna](https://github.com/mna)).

This package is currently under development, no API guaranties.

## Install

```
go get -t github.com/breml/logstash-config/...
```

## Usage

### ls-config-check

`cmd/ls-config-check` contains the source code for a small sample tool, which uses just allow to parse a given Logstash config file. If the file could be parsed successfully, the tool just exits with exit code `0`. If the parsing fails, the exit code is non zero and a error message, indicating the location, where the parsing failed, is printed.
`ls-config-check <logstash-config-file>` could be used instead if `bin/logstash -f <logstash-config-file> -t`, which is orders of magnitude faster 😃.

## Rebuild parser

1. Get and install [pigeon](https://github.com/mna/pigeon).
2. Run `go generate` in the root directory of this repository.

## Author

Copyright 2017 by Lucas Bremgartner ([breml](https://github.com/breml))

## License

[Apache 2.0](LICENSE)
