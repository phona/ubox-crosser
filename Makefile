SOURCES := $(shell find . -name "*.go")
BIN_PATH := ../../bin
export GOOS=linux
export GOARCH=amd64

all: $(SOURCES)
	@echo "Building $@ ..."
	cd cmd/client; go install
	@echo "Building $@ ..."
	cd cmd/server; go install
	
clean: 
	rm -f $(BIN_PATH)/linux_amd64/client $(BIN_PATH)/linux_amd64/server

client: $(SOURCES) cmd/client/client.go
	@echo "Building $@ ..."
	cd cmd/client; go install

server: $(SOURCES) cmd/server/server.go
	@echo "Building $@ ..."
	cd cmd/server; go install
