all: main.go fetcher.go parser.go storage.go url_manager.go info.go
	mkdir -p output
	time go build -o crawler *.go
	# time ./crawler -cpuprofile=cpu.prof -memprofile=mem.prof 
	time ./crawler  

clean: 
	time go clean
