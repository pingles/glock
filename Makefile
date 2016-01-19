.PHONY: clean

GOPATH:=${shell pwd}

croney:
	GO15VENDOREXPERIMENT=1 go install croney
clean:
	rm -rf ./bin
	rm -rf ./pkg
