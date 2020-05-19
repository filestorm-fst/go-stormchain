# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: moac vnode all test clean
.PHONY: moac-linux moac-linux-386 moac-linux-amd64 moac-linux-mips64 moac-linux-mips64le
.PHONY: moac-linux-arm moac-linux-arm-5 moac-linux-arm-6 moac-linux-arm-7 moac-linux-arm64
.PHONY: moac-darwin moac-darwin-386 moac-darwin-amd64
.PHONY: moac-windows moac-windows-386 moac-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

compile_chain3:
	cd internal/jsre/deps && go-bindata -nometadata -pkg deps -o bindata.go bignumber.js chain3.js && gofmt -w -s bindata.go;

#build/env.sh go run build/mci.go install ./cmd/moac
moac: compile_chain3
	build/env.sh go run build/ci.go install ./cmd/moac

	@echo "Done building."
	@echo "Run \"$(GOBIN)/moac\" to launch moac."

vnode: compile_chain3
	build/env.sh go run build/ci.go install ./cmd/moac
    #go run build/ci.go install ./cmd/moac
	@echo "Done building Vnode."
	@echo "Run \"$(GOBIN)/moac\" to launch moac."
#Note, this not include SCS yet, 2017/11/22
all:
	build/env.sh go run build/mci.go install

android:
	build/env.sh go run build/mci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/moac.aar\" to use the library."

bootnode:
	build/env.sh go build -o ./build/bin/bootnode cmd/bootnode/main.go

ios:
	build/env.sh go run build/mci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Moac.framework\" to use the library."

test: all
	build/env.sh go run build/mci.go test

clean:
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

reset:
	rm -fr build/_workspace/pkg/ $(GOBIN)/* build/_workspace/src/

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/jteeuwen/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go install ./cmd/abigen

# Cross Compilation Targets (xgo)

moac-cross: moac-linux moac-darwin moac-windows moac-android moac-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/moac-*

moac-linux: moac-linux-386 moac-linux-amd64 moac-linux-arm moac-linux-mips64 moac-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/moac-linux-*

moac-linux-386:
	build/env.sh go run build/mci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/moac
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/moac-linux-* | grep 386

moac-linux-amd64: compile_chain3
	build/env.sh go run build/mci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/moac
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/moac-linux-* | grep amd64

moac-linux-arm: moac-linux-arm-5 moac-linux-arm-6 moac-linux-arm-7 moac-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/moac-linux-* | grep arm

moac-linux-arm-5:
	build/env.sh go run build/mci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/moac
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/moac-linux-* | grep arm-5

moac-linux-arm-6:
	build/env.sh go run build/mci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/moac
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/moac-linux-* | grep arm-6

moac-linux-arm-7:
	build/env.sh go run build/mci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/moac
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/moac-linux-* | grep arm-7

moac-linux-arm64:
	build/env.sh go run build/mci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/moac
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/moac-linux-* | grep arm64

moac-linux-mips:
	build/env.sh go run build/mci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/moac
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/moac-linux-* | grep mips

moac-linux-mipsle:
	build/env.sh go run build/mci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/moac
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/moac-linux-* | grep mipsle

moac-linux-mips64:
	build/env.sh go run build/mci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/moac
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/moac-linux-* | grep mips64

moac-linux-mips64le:
	build/env.sh go run build/mci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/moac
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/moac-linux-* | grep mips64le

moac-darwin: moac-darwin-386 moac-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/moac-darwin-*

moac-darwin-386:
	build/env.sh go run build/mci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/moac
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/moac-darwin-* | grep 386

moac-darwin-amd64: compile_chain3
	build/env.sh go run build/mci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/moac
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/moac-darwin-* | grep amd64

moac-windows: moac-windows-386 moac-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/moac-windows-*

moac-windows-386:
	build/env.sh go run build/mci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/moac
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/moac-windows-* | grep 386

moac-windows-amd64:
	build/env.sh go run build/mci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/moac
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/moac-windows-* | grep amd64
