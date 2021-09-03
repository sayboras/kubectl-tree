BINARY ?= dist/kubectl-tree
INSTALL_LOCATION ?= ~/.krew/store/tree/v0.4.0/

build:
	go build -o $(BINARY) ./cmd/kubectl-tree

install:
	cp $(BINARY) $(INSTALL_LOCATION)