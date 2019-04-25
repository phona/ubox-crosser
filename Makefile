SOURCES := $(shell find . -name "*.go")
BIN_PATH := ../../bin
export GOOS=linux
export GOARCH=amd64

all: client server

clean: 
	rm -f $(BIN_PATH)/linux_amd64/client $(BIN_PATH)/linux_amd64/server

client: $(SOURCES)
	@echo "Building $@ ..."
	cd cmd/client; go install

server: $(SOURCES)
	@echo "Building $@ ..."
	cd cmd/server; go install
