# logstash-config : parser and abstract syntax tree for [Logstash](https://www.elastic.co/logstash/) config files

[![Test Status](https://github.com/breml/logstash-config/workflows/Test/badge.svg)](https://github.com/breml/logstash-config/actions?query=workflow%3ATest)
 [![Go Report Card](https://goreportcard.com/badge/github.com/breml/logstash-config)](https://goreportcard.com/report/github.com/breml/logstash-config)\
[![GoDoc](https://pkg.go.dev/badge/github.com/breml/logstash-config)](https://pkg.go.dev/github.com/breml/logstash-config) [![License](https://img.shields.io/badge/license-Apache_2.0-blue.svg)](LICENSE)

## Overview

The Go package config provides a ready to use parser for Logstash ([github](https://github.com/elastic/logstash)) configuration files.

The basis of the grammar for the parsing of the Logstash configuration format is the original [Logstash Treetop grammar](https://github.com/elastic/logstash/blob/master/logstash-core/lib/logstash/config/grammar.treetop) which could be used with only minor changes.

logstash-config uses [pigeon](https://github.com/mna/pigeon) to generate the parser from the PEG (parser expression grammar). Special thanks to Martin Angers ([mna](https://github.com/mna)).

This package is currently under development, no API guaranties.

## Install

```bash
go get -t github.com/breml/logstash-config/...
```

## Usage

### mustache

`mustache` is a command line tool that allows to syntax check, lint and format Logstash configuration files. The name of
the tool is inspired by the original Logstash Logo ([wooden character with an eye-catching mustache](https://www.elastic.co/de/blog/high-level-logstash-roadmap-is-published)).

The `check` command verifies the syntax of Logstash configuration files:

```shell
mustache check file.conf
```

The `lint` command checks for problems in Logstash configuration files.

The following checks are performed:

* Valid Logstash configuration file syntax
* No comments in exceptional places (these are comments, that are valid by the Logstash configuration file syntax, but
  but are located in exceptional or uncommon locations)
* Precence of an `id` attribute for each plugin in the Logstash configuration

If the `--auto-fix-id` flag is passed, each plugin gets automatically an ID. Be aware, that this potentially reformats
the Logstash configuration files.

```shell
mustache lint --auto-fix-id file.conf
```

With the `format` command, mustache returns the provided configuration files in a standardized format (indentation,
location of comments). By default, the reformatted file is print to standard out. If the flag `--write-to-source`
is provided, the Logstash config files are reformatted in place.

```shell
mustache format --write-to-source file.conf
```

Use the `--help` flag to get more information about the usage of the tool.

## Rebuild parser

1. Get and install [pigeon](https://github.com/mna/pigeon).
2. Run `go generate` in the root directory of this repository.

## Author

Copyright 2017-2021 by Lucas Bremgartner ([breml](https://github.com/breml))

## License

[Apache 2.0](LICENSE)
