
PREFIX ?= /usr/local

DIST := COPYING Makefile readme.md

PROJECT := mtmapprune
VERSION = 3
BUILD = `git describe --tags --always`

$(PROJECT): main.go
	go build -ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)" -o $(PROJECT)

build: $(PROJECT)

install: $(PROJECT)
	go install
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	install -m0755 $(PROJECT) $(DESTDIR)$(PREFIX)/bin/$(PROJECT)

clean:
	go clean

dist:
	rm -rf $(PROJECT)-$(BUILD)
	mkdir $(PROJECT)-$(BUILD)
	cp $(DIST) $(PROJECT)-$(BUILD)/
	rm -f $(PROJECT)-$(BUILD)/$(PROJECT)
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)" -o $(PROJECT) -o $(PROJECT)-$(BUILD)/$(PROJECT)
	rm -f $(PROJECT)-$(BUILD)-x86_64.zip
	zip -r $(PROJECT)-$(BUILD)-x86_64.zip $(PROJECT)-$(BUILD)/
	rm -f $(PROJECT)-$(BUILD)/$(PROJECT)
	CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)" -o $(PROJECT) -o $(PROJECT)-$(BUILD)/$(PROJECT).exe
	rm -f $(PROJECT)-$(BUILD)-win64.zip
	zip -r $(PROJECT)-$(BUILD)-win64.zip $(PROJECT)-$(BUILD)/

