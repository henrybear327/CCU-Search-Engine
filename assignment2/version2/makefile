all: chromedp.go crawler.go fetcher.go init.go main.go manager.go output.go parser.go preprocess.go storage.go
	# mkdir -p output
	time go build -o crawler *.go
	# time ./crawler -cpuprofile=cpu.prof -memprofile=mem.prof 
	time ./crawler  

clean: 
	time go clean
