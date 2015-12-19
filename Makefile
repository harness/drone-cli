BIN  := bin
DIST := dist

SRC = $(wildcard drone/*.go)
RELEASES = $(DIST)/drone_linux_amd64.tar.gz \
	   $(DIST)/drone_linux_386.tar.gz \
	   $(DIST)/drone_linux_arm.tar.gz \
	   $(DIST)/drone_darwin_amd64.tar.gz \
	   $(DIST)/drone_windows_386.tar.gz \
	   $(DIST)/drone_windows_amd64.tar.gz

GO = GO15VENDOREXPERIMENT=1 go

install: $(BIN)/drone
	cp $< $(GOPATH)/bin/

release: $(RELEASES)

deps:
	go get -u ./drone/...

test:
	$(GO) test ./drone/...

clean:
	rm -rf $(BIN) $(DIST)

$(BIN)/drone: $(SRC)
	$(GO) build -o $@ $(SRC)

.PRECIOUS: $(BIN)/%/drone
$(BIN)/%/drone: GOOS=$(firstword $(subst _, ,$*))
$(BIN)/%/drone: GOARCH=$(subst .exe,,$(word 2,$(subst _, ,$*)))
$(BIN)/%/drone: $(SRC)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $@ $(SRC)

$(DIST)/drone_%.tar.gz: $(BIN)/%/drone
	mkdir -p $(DIST)
	tar -cvzf $@ --directory=$(BIN)/$* drone
	sha256sum $@ > $(DIST)/drone_$*.sha256
