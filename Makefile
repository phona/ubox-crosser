export GOOS=linux
export GOARCH=amd64
SOURCES := $(shell find . -name "*.go")

all: 
	$(client) $(server)

clean: 
	rm -f $(GOPATH)/bin/client $(GOPATH)/bin/server

client: $(SOURCES) 
	@echo "Building $@ ..."
	cd cmd/client; go install

server: $(SOURCES)
	@echo "Building $@ ..."
	cd cmd/server; go install
