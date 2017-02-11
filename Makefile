NAME := docker-ipv6nat
PKG := github.com/robbertkl/$(NAME)
DIR := /go/src/$(PKG)
GO := 1.7.5-alpine3.5
TAG := `git describe --tags`
LDFLAGS := -X main.buildVersion=$(TAG)
TARGETS := $(NAME) $(NAME).aarch64 $(NAME).armhf

.SILENT:

.PHONY: all
all: clean $(TARGETS)

.PHONY: clean
clean:
	rm -f $(TARGETS)

$(NAME):
	docker run --rm \
		-v "$(PWD)":"$(DIR)" \
		-w "$(DIR)" \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-e CGO_ENABLED=0 \
		golang:"$(GO)" \
		go build -o "$@" -ldflags "$(LDFLAGS)" "./cmd/$(NAME)"

$(NAME).aarch64:
	docker run --rm \
		-v "$(PWD)":"$(DIR)" \
		-w "$(DIR)" \
		-e GOOS=linux \
		-e GOARCH=arm64 \
		-e CGO_ENABLED=0 \
		golang:"$(GO)" \
		go build -o "$@" -ldflags "$(LDFLAGS)" "./cmd/$(NAME)"

$(NAME).armhf:
	docker run --rm \
		-v "$(PWD)":"$(DIR)" \
		-w "$(DIR)" \
		-e GOOS=linux \
		-e GOARCH=arm \
		-e GOARM=7 \
		-e CGO_ENABLED=0 \
		golang:"$(GO)" \
		go build -o "$@" -ldflags "$(LDFLAGS)" "./cmd/$(NAME)"
