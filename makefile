# Go commands
GOCMD=go
RICECMD=rice
RSRCCMD=rsrc
GOBUILD=$(GOCMD) build
BINARY_NAME=HorribleNotifier

all: build
ifeq ($(OS),Windows_NT)
build::
	$(RICECMD) embed-go
	$(RSRCCMD) -manifest horriblenotifier.manifest -ico icons\hn-128.ico -o rsrc.syso
	$(GOBUILD) -ldflags "-H windowsgui -s -w" -o $(BINARY_NAME).exe
else
build::
	$(GOBUILD) -o $(BINARY_NAME)
	strip $(BINARY_NAME)
endif