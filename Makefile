help:
	@@echo make lint
	@@echo make build
	@@echo make run
	@@echo make clean

lint:
	golangci-lint run

build: main.go
	go build -o hello-concur-fileio-go

run: build
	go run hello-concur-fileio-go

clean:
	rm hello-concur-fileio-go ints.txt
