build:
	godep go build .

clean:
	-rm *.log
	-rm jmcom
	-rm jmcom_linux*
	-rm jmcom_windows*
	-rm jmcom_darwin*

build_all:
	@echo "Setting GOPATHS for build"
	@export GOPATHVAR=$(GOPATH)
	@export GOPATHDEP=$(shell godep path)
	@export GOPATH=$(GOPATHDEP)
	@echo "Building for linux"
	gox -osarch="linux/amd64"
	gox -osarch="linux/386"
	gox -osarch="linux/arm"
	@echo "Building for Windows"
	gox -osarch="windows/amd64"
	gox -osarch="windows/386"
	@echo "Building for Mac OSX"
	gox -osarch="darwin/amd64"
	gox -osarch="darwin/386"
	@export GOPATH=$(GOPATHVAR)
