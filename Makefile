.PHONY: download calc

all: download calc
download:
	go build -o bin/download code/download/*.go
calc:
	go build -o bin/calc code/calc/*.go