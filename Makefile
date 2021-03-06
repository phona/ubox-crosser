SOURCES := $(shell find . -name "*.go")
BIN_PATH := ../../bin
export GOOS=linux
export GOARCH=amd64

all: $(SOURCES)
	@echo "Building $@ ..."
	cd cmd/client; go install
	@echo "Building $@ ..."
	cd cmd/server; go install
	@echo "Building $@ ..."
	cd cmd/auth_server; go install

clean: 
	rm -f $(BIN_PATH)/linux_amd64/client $(BIN_PATH)/linux_amd64/server
