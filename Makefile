NAME=chatgpt
BINDIR=bin
GOFILES=./cmd
VERSION=$(shell git describe --tags || echo "%version")
BUILDTIME=$(shell date -u)
GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-X "github.com/vvtommy/chatgpt-cli/cmd.Version=$(VERSION)" \
		-X "github.com/vvtommy/chatgpt-cli/cmd.BuildTime=$(BUILDTIME)" \
		-w -s -buildid='

PLATFORM_LIST = \
	darwin-amd64 \
	darwin-amd64-v3 \
	darwin-arm64 \
	linux-386 \
	linux-amd64 \
	linux-amd64-v3 \
	linux-armv5 \
	linux-armv6 \
	linux-armv7 \
	linux-arm64 \
	linux-mips-softfloat \
	linux-mips-hardfloat \
	linux-mipsle-softfloat \
	linux-mipsle-hardfloat \
	linux-mips64 \
	linux-mips64le \
	linux-riscv64

WINDOWS_ARCH_LIST = \
	windows-386 \
	windows-amd64 \
	windows-amd64-v3 \
	windows-arm64 \
	windows-armv7

all: linux-amd64 darwin-amd64 windows-amd64 # Most used

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

darwin-amd64-v3:
	GOARCH=amd64 GOOS=darwin GOAMD64=v3 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

darwin-arm64:
	GOARCH=arm64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

linux-386:
	GOARCH=386 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

linux-amd64-v3:
	GOARCH=amd64 GOOS=linux GOAMD64=v3 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

linux-armv5:
	GOARCH=arm GOOS=linux GOARM=5 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

linux-armv6:
	GOARCH=arm GOOS=linux GOARM=6 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

linux-armv7:
	GOARCH=arm GOOS=linux GOARM=7 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

linux-arm64:
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

linux-mips-softfloat:
	GOARCH=mips GOMIPS=softfloat GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

linux-mips-hardfloat:
	GOARCH=mips GOMIPS=hardfloat GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

linux-mipsle-softfloat:
	GOARCH=mipsle GOMIPS=softfloat GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

linux-mipsle-hardfloat:
	GOARCH=mipsle GOMIPS=hardfloat GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

linux-mips64:
	GOARCH=mips64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

linux-mips64le:
	GOARCH=mips64le GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

linux-riscv64:
	GOARCH=riscv64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(GOFILES)

windows-386:
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe $(GOFILES)

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe $(GOFILES)

windows-amd64-v3:
	GOARCH=amd64 GOOS=windows GOAMD64=v3 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe $(GOFILES)

windows-arm64:
	GOARCH=arm64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe $(GOFILES)

windows-armv7:
	GOARCH=arm GOOS=windows GOARM=7 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe $(GOFILES)

gz_releases=$(addsuffix .gz, $(PLATFORM_LIST))
zip_releases=$(addsuffix .zip, $(WINDOWS_ARCH_LIST))

$(gz_releases): %.gz : %
	chmod +x $(BINDIR)/$(NAME)-$(basename $@)
	gzip -f -S -$(VERSION).gz $(BINDIR)/$(NAME)-$(basename $@)

$(zip_releases): %.zip : %
	zip -m -j $(BINDIR)/$(NAME)-$(basename $@)-$(VERSION).zip $(BINDIR)/$(NAME)-$(basename $@).exe

all-arch: $(PLATFORM_LIST) $(WINDOWS_ARCH_LIST)

releases: $(gz_releases) $(zip_releases)

lint:
	GOOS=darwin golangci-lint run ./...
	GOOS=windows golangci-lint run ./...
	GOOS=linux golangci-lint run ./...

clean:
	rm $(BINDIR)/*

