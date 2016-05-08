GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
RELEASE_DIR=releases
SRC_FILES = $(wildcard *.go)

build-windows-amd64:
	@$(MAKE) build GOOS=windows GOARCH=amd64 SUFFIX=.exe

build-linux-amd64:
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-darwin-amd64:
    @$(MAKE) build GOOS=darwin GOARCH=amd64
    
build-linux-arm5:
    @$(MAKE) build GOOS=linux GOARCH=arm GOARM=5

$(RELEASE_DIR)/$(GOOS)/$(GOARCH)/zwaymqtt$(SUFFIX): $(SRC_FILES)
	go build -o $(RELEASE_DIR)/$(GOOS)/$(GOARCH)/zwaymqtt$(SUFFIX) .

build: $(RELEASE_DIR)/$(GOOS)/$(GOARCH)/zwaymqtt$(SUFFIX)

all:
	@$(MAKE) build-windows-amd64 
	@$(MAKE) build-linux-amd64
	@$(MAKE) build-darwin-amd64
	@$(MAKE) build-linux-arm5