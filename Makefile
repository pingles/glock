.PHONY: clean

GOPATH:=${shell pwd}

./bin/glock: ${wildcard src/glock/**/*.go}
	GO15VENDOREXPERIMENT=1 go install glock

clean:
	rm -rf ./bin
	rm -rf ./pkg
