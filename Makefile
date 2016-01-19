.PHONY: clean

croney:
	GO15VENDOREXPERIMENT=1 go build -o croney .
clean:
	rm -f ./croney
