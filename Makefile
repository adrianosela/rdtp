NAME:=$(shell basename `git rev-parse --show-toplevel`)
HASH:=$(shell git rev-parse --verify --short HEAD)

all: build

build:
	go build -ldflags "-X main.version=$(HASH)" -o $(NAME)
