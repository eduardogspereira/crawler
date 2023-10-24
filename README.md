<img src="https://avatars.githubusercontent.com/u/11029372?s=200&v=4" width="127px" height="127px" align="left"/>

# Crawler

Crawler built as part of Monzo Take Home Test.

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