SRC := cmd/paletter/main.go $(wildcard generator/*) colors.go doc.go paletter.go
LDFLAGS := -s -w

all: build/linux.tar build/mac.tar build/windows.tar

builddir:
	@mkdir -p build/linux build/windows build/mac

build/linux/paletter: $(SRC) builddir
	@echo -e "\n# Building for Linux"
	env GOOS=linux GOARCh=amd64 go build -ldflags="-s -w" -o $@ ./cmd/paletter

build/linux.tar: build/linux/paletter
	tar cf $@ -C $(dir $@)linux $(notdir $(wildcard build/linux/*))

build/mac/paletter.app: $(SRC) dist/Info.plist builddir
	@# See https://developer.apple.com/library/archive/documentation/CoreFoundation/Conceptual/CFBundles/BundleTypes/BundleTypes.html
	@# and https://apple.stackexchange.com/questions/253184/associating-protocol-handler-in-mac-os-x
	@echo -e "\n# Building MacOS app bundle"
	@mkdir -p $@/Contents
	env GOOS=darwin GOARCh=amd64 go build -ldflags="-s -w" -o $@/Contents/MacOS/paletter ./cmd/paletter
	cp dist/Info.plist $@/Contents

build/mac.tar: build/mac/paletter.app
	tar cf $@ -C $(dir $@)/mac paletter.app

build/windows/paletter.exe: $(SRC) builddir
	@echo -e "\n# Building for Windows"
	env GOOS=windows GOARCh=amd64 go build -ldflags="-s -w -H windowsgui" -o $@ ./cmd/paletter

build/windows.tar: build/windows/paletter.exe
	tar cf $@ -C $(dir $@)windows $(notdir $(wildcard build/windows/*))

clean:
	rm -rf build/*

test:
	go test -v -benchmem ./...

.PHONY: all builddir clean test