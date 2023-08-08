LINKERFLAGS = -X main.Version=`git describe --tags --always --dirty` -X main.BuildTimestamp=`date -u '+%Y-%m-%d_%I:%M:%S_UTC'`
PROJECTROOT = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
PROJECTNAME = $(lastword $(subst /, ,$(PROJECTROOT)))


all: clean build

.PHONY: clean release
clean:
	echo $(PROJECTNAME)
	@echo Running clean job...
	rm -f coverage.txt
	rm -rf bin/ release/
	rm -f main busydonkey



build:
	@echo Running build job...
	mkdir -p bin/linux/arm bin/linux/x64 bin/windows bin/osx/x64 bin/osx/arm
	GOOS=linux GOARCH=arm64 go build  -ldflags "$(LINKERFLAGS)" -o bin/linux/arm ./...
	GOOS=linux GOARCH=amd64 go build  -ldflags "$(LINKERFLAGS)" -o bin/linux/x64 ./...
	GOOS=windows GOARCH=amd64 go build  -ldflags "$(LINKERFLAGS)" -o bin/windows ./...
	GOOS=darwin GOARCH=amd64 go build  -ldflags "$(LINKERFLAGS)" -o bin/osx/x64 ./...
	GOOS=darwin GOARCH=arm64 go build  -ldflags "$(LINKERFLAGS)" -o bin/osx/arm ./...


release: build
	mkdir -p release
	$(eval VER=$(shell sh -c "bin/osx/x64/$(PROJECTNAME) -version |cut -f 2 -d ' '"))
	cd bin && tar -zcpv -s /linux/$(PROJECTNAME)-linux-$(VER)/ -f ../release/$(PROJECTNAME)-linux-$(VER).tgz linux/*
	cd bin/windows &&  zip -r -9 ../../release/$(PROJECTNAME)-win-$(VER).zip *

