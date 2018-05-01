# Version 2

Be polite, concurrent, and seed from [Alexa](https://www.alexa.com/topsites)

* Master (Scheduler)
* Fetcher (Downloader)
* Parser
    * robots.txt, sitemap.xml
    * general webpages
* URL Manager (queuing system)

# Golang

## external packages

* goquery: html parser
* toml: config file
* chromedp: headless browser
* logrus: logger
* go-sqlite3: sqlite3 driver
* goOse: main text extraction
* Oneliner
```bash
go get github.com/PuerkitoBio/goquery
go get github.com/BurntSushi/toml
```

## internal packages

* log
* http
* url
* publicsuffix

# Notes

## [`[]byte` and `string` conversion](https://studygolang.com/articles/10526)

* string to []byte
```go
var str string = "test"
var data []byte = []byte(str)
```

* []byte to string
```go
var data [10]byte 
byte[0] = 'T'
byte[1] = 'E'
var str string = string(data[:])
```

## `[]byte` to `io.Reader`

```go
res := bytes.NewReader(pageSource)
```

## Function timer

```go
start := time.Now()

elapsed := time.Since(start)
log.Printf("GetPageSource took %s", elapsed)
```

## Profiler

* [Reading](https://golang.org/pkg/runtime/pprof/)
* `brew install Graphviz`
* `go tool pprof cpu.prof` and type web for more detail