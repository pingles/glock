.PHONY: clean

GOPATH:=${shell pwd}
VERSION?=dev

./bin/glock: ${wildcard src/glock/*.go}
	GO15VENDOREXPERIMENT=1 go install -ldflags "-X main.version=${VERSION}" glock

clean:
	rm -rf ./bin
	rm -rf ./pkg
