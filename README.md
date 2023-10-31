# Crawler

Playing around with Go.

### Dependencies

- [Go ^1.21.3](https://golang.org/dl/)

### Building the crawler

```shell
go build -o crawler .
```

### Running the crawler
```shell
./crawler --help                                                                                                                                                                        00:42:51
Usage of ./crawler:
  -r, --retries uint   Number of task retries (default 3)
  -t, --timeout int    HTTP timeout (seconds) (default 30)
  -u, --url string     Target URL
  -w, --workers int    Number of workers (default 100)
pflag: help requested
```

### Running the tests

```shell
go test -v .
```
