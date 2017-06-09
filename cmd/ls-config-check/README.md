# ls-config-check

`cmd/ls-config-check` contains the source code for a small sample tool, which uses just allow to parse a given Logstash config file. If the file could be parsed successfully, the tool just exits with exit code `0`. If the parsing fails, the exit code is non zero and a error message, indicating the location, where the parsing failed, is printed.
`ls-config-check <logstash-config-file>` could be used instead if `bin/logstash -f <logstash-config-file> -t`, which is orders of magnitude faster 😃. 

## Install

```
go get -t github.com/breml/logstash-config/cmd/ls-config-check
```

## Author

Copyright 2017 by Lucas Bremgartner ([breml](https://github.com/breml))

## License

[Apache 2.0](LICENSE)