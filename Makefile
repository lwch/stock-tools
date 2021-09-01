.PHONY: download calc list filter

all: download calc list filter
download:
	go build -o bin/download code/download/*.go
calc:
	go build -o bin/calc code/calc/*.go
list:
	go build -o bin/list code/list/*.go
filter:
	go build -o bin/filter code/filter/*.go
clean:
	rm -fr bin