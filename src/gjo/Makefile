GOBUILD=go build
GOCLEAN=go clean
BINARY_NAME=gjo

export GO11MODULE=on

all: build

clean:
	$(GOCLEAN)

build: clean
	$(GOBUILD) -o $(BINARY_NAME)

# copy to $GOBIN
install: build
	cp -f $(BINARY_NAME) $(GOBIN)/

# build release binary
release: clean
	GOOS=darwin GOARCH=amd64 $(GOBUILD) && zip MacOS.zip $(BINARY_NAME) && rm -rf $(BINARY_NAME)
	GOOS=linux GOARCH=amd64 $(GOBUILD) && zip Linux.zip $(BINARY_NAME) && rm -rf $(BINARY_NAME)
	GOOS=windows GOARCH=amd64 $(GOBUILD) && zip Windows.zip $(BINARY_NAME).exe && rm -rf $(BINARY_NAME).exe


