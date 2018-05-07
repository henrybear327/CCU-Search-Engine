# Version 2

> Politeness Performance Scalability 

Seed sites from [Alexa](https://www.alexa.com/topsites)

* Master (Scheduler)
* Fetcher (Downloader)
* Parser
    * robots.txt, sitemap.xml
    * general webpages
* URL Manager (queuing system)

[slide](https://docs.google.com/presentation/d/107ohaNKpRw_JEpvXizm_QLTbT3M12KQYpf_dYtLFDBE/edit?usp=sharing)

# Golang

## Development notes

1. map is not thread safe, use go channels to pass data back and forth
2. robots might be missing in robotsCollection for some url -> no robots.txt -> no limitations
3. Trim string! Can't emphasize more!
4. `google.com.tw` tld != `google.com`
5. Goroutines still comes with cost! Think twice before creating them
6. Watch out for deadlock!
    * Acquire locks all following the same sequence
7. `pkill "Google Chrome"` you will need this!

## external packages

* goquery: html parser
* BurntSushi/toml: config file
* robotstxt: parse robots.txt rules
* color: terminal output color
* chromedp: headless browser
    * mac must install `brew install Caskroom/versions/google-chrome-canary`
* globalsign/mgo: mongodb driver
* RadhiFadlillah/go-readability: main text extraction

### Installation

```bash
go get github.com/PuerkitoBio/goquery
go get github.com/BurntSushi/toml
go get github.com/temoto/robotstxt
go get github.com/fatih/color
go get -u github.com/chromedp/chromedp
go get github.com/globalsign/mgo
go get github.com/henrybear327/go-readability
```

## internal packages

* log: print messages like system logs with ease
* http: static site fetching
* url: decompose the link into hostname + path + query ...
* publicsuffix: find out the tld for any given link (determine if the link is a internal/external link)

## ROBO3T

MongoDB GUI `brew cask install robo-3t`

# Virtual env

* `sudo pip3 install virtualenv `
* install
```bash
pip install Flask elasticsearch pymongo
```

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

## `io.Reader` to `[]byte`

```go
robots, err := ioutil.ReadAll(res.Body)
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
fmt.Println("get host name", u.Hostname())
```

## XML parsing 

* [reading](https://tutorialedge.net/golang/parsing-xml-with-golang/)
* [what is XML](https://www.awoo.com.tw/blog/2018/01/sitemap-xml/)

## Map locking

* [reading](https://blog.golang.org/go-maps-in-action)
    * `RWLock` has `Lock` and `Rlock`

## defer must be a func

There you go

```go
defer func(done chan bool) {
    done <- true
}(done)
```

## SQLite3

* [reading](https://astaxie.gitbooks.io/build-web-application-with-golang/zh/05.3.html)

## Time

* to string
```go
fetchTime.Format(time.RFC3339)
```
* from string
```go
fetchTime, e := time.Parse(
    time.RFC3339,
    "2012-11-01T22:08:41+00:00")
p(fetchTime)
```