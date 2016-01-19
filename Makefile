.PHONY: clean

GOPATH:=${shell pwd}

./bin/crony: ${wildcard src/**/*.go}
	GO15VENDOREXPERIMENT=1 go install crony
clean:
	rm -rf ./bin
	rm -rf ./pkg
run: ./bin/crony
	./bin/crony --lockPath=/crony/testing --schedule="* * * * *" --command="/usr/bin/env pwd"