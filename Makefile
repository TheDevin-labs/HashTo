BINARY := hashto
PREFIX ?= /usr/local

.PHONY: build install clean

build:
	go build -o $(BINARY) .

install: build
	install -Dm755 $(BINARY) $(PREFIX)/bin/$(BINARY)

clean:
	rm -f $(BINARY) *.hash *.json *.yaml
