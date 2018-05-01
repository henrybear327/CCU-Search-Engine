all: main.go fetcher.go parser.go storage.go url_manager.go info.go
	time go build -o crawler *.go
	time ./crawler
	
clean: 
	time go clean
