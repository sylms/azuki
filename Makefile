.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	go build -o ./azuki

.PHONY: run
run: build
	./azuki
