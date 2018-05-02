# Version 2

Be polite, concurrent, and seed from [Alexa](https://www.alexa.com/topsites)

* Master (Scheduler)
* Fetcher (Downloader)
* Parser
    * robots.txt, sitemap.xml
    * general webpages
* URL Manager (queuing system)

# Golang

## Development notes

1. map is not thread safe, use go channels to pass data back and forth
2. robots might be missing in robotsCollection for some url -> no robots.txt -> no limitations

## external packages

* goquery: html parser
* BurntSushi/toml: config file
* robotstxt: parse robots.txt rules
* color: terminal output color

### candidates 

* chromedp: headless browser
* logrus: logger
* go-sqlite3: sqlite3 driver
* goOse: main text extraction
* Oneliner

### Installation

```bash
go get github.com/PuerkitoBio/goquery
go get github.com/BurntSushi/toml
go get github.com/temoto/robotstxt
go get github.com/fatih/color
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

## Wait for all goroutines to finish

* [reading](https://nathanleclaire.com/blog/2014/02/15/how-to-wait-for-all-goroutines-to-finish-executing-before-continuing/)

## URL decomposition

```go
u, err := url.Parse(link)
if err != nil {
	log.Fatal(err)
}
fmt.Println("url decompose", u.Scheme, u.Host, u.Path, u.RawQuery)
```