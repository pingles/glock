.PHONY: clean

GOPATH:=${shell pwd}

crony: ${wildcard src/crony/**/*.go}
	GO15VENDOREXPERIMENT=1 go install crony
clean:
	rm -rf ./bin
	rm -rf ./pkg
