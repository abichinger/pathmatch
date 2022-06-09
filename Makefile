test: 
	go test -race ./...

bench:
	go test -benchmem -bench=.

lint:
	golangci-lint run --verbose