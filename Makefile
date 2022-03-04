GOCMD=go
GOPATH=$(shell $(GOCMD) env GOPATH))
GOTEST=$(GOCMD) test

all: test

test:
	$(GOTEST) -v ./...  -cover
	(cd ./gousujwt/ && $(GOTEST) -v ./... -cover)
	(cd ./gousuredis/ && $(GOTEST) -v ./... -cover)
	(cd ./gousuchi/ && $(GOTEST) -v ./... -cover)
	(cd ./gousupostgres/ && $(GOTEST) -v ./... -cover)
	(cd ./goususmtp/ && $(GOTEST) -v ./... -cover)
	(cd ./gousuldap/ && $(GOTEST) -v ./... -cover)
	(cd ./gousukafka/ && $(GOTEST) -v ./... -cover)
	(cd ./gousustomp/ && $(GOTEST) -v ./... -cover)
