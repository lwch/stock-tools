.PHONY: download info list filter dnn

all: download info list filter dnn
download:
	go build -o bin/download code/download/*.go
info:
	go build -o bin/info code/info/*.go
list:
	go build -o bin/list code/list/*.go
filter:
	go build -o bin/filter code/filter/*.go
dnn:
	go build -o bin/dnn code/dnn/*.go
clean:
	rm -fr bin