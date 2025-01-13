.PHONY: all install uninstall clean dvi dist tests

all: build

build:
	go build -o maze main.go

install:
	mv maze /usr/local/bin/

uninstall:
	rm /usr/local/bin/maze

clean:
	go clean

dvi:
	echo "No DVI generation in this project."

dist:
	echo "No distribution in this project."

tests:
	go test ./...
