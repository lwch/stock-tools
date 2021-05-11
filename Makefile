.PHONY: download calc list

all: download calc list
download:
	go build -o bin/download code/download/*.go
calc:
	go build -o bin/calc code/calc/*.go
list:
	go build -o bin/list code/list/*.go
clean:
	rm -fr bin