all: main.go init.go fetcher.go parser.go storage.go manager.go output.go preprocess.go chromedp.go crawler.go
	mkdir -p output
	time go build -o crawler *.go
	# time ./crawler -cpuprofile=cpu.prof -memprofile=mem.prof 
	time ./crawler  

clean: 
	time go clean
