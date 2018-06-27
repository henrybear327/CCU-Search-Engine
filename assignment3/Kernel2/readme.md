# Search Kernel

Simple inverted page searching kernel

[Presentation slides](https://docs.google.com/presentation/d/1YRgBzzz5Y6f5qyWeQRvchEjwKv-QTGtD_-AugWRorFE/edit?usp=sharing)

## Usage

* Insert
```bash
curl -X POST \
-d '{"title":"title", "body":"hello world", "url":"google.com"}' \
http://localhost:8001/insert
```
* Search
```bash
curl -X POST \
-d '{"query":"ettoday"}' \
http://localhost:8001/search
```

## Package

* gse
    * segmentation, `go get -u github.com/go-ego/gse`
* gopsutil
    * system stat, `go get -u github.com/shirou/gopsutil/mem`
    * code
```go
v, _ := mem.VirtualMemory()
    megabyte := uint64(1024 * 1024)
    // almost every return value is a struct
    log.Printf("Total: %v MB, Free:%v MB, UsedPercent:%f%%\n", v.Total/megabyte, v.Free/megabyte, v.UsedPercent)
    // convert to JSON. String() is also implemented
    // log.Println(v)
```