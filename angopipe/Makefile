.PHONY: build clean

build:
	go build -o build/ango cmd/ango/main.go
	go build -o build/fukugo cmd/fukugo/main.go
	go build -o build/keygen cmd/keygen/main.go

clean:
	rm -f build/ango
	rm -f build/fukugo
	rm -f build/keygen
