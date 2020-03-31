NAME:=$(shell basename `git rev-parse --show-toplevel`)
RELEASE:=$(shell git rev-parse --verify --short HEAD)
VERSION = 0.1.0

all: setbin

clean:
	rm -f $(NAME)

setbin: build
	cp $(NAME) /usr/local/bin

build:
	go build -ldflags "-X main.version=$(VERSION)-$(RELEASE)" -o $(NAME) ./cli
