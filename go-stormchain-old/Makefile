# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: storm android ios storm-cross evm all test clean
.PHONY: storm-linux storm-linux-386 storm-linux-amd64 storm-linux-mips64 storm-linux-mips64le
.PHONY: storm-linux-arm storm-linux-arm-5 storm-linux-arm-6 storm-linux-arm-7 storm-linux-arm64
.PHONY: storm-darwin storm-darwin-386 storm-darwin-amd64
.PHONY: storm-windows storm-windows-386 storm-windows-amd64

GOBIN = ./build/bin
GO ?= latest

storm:
	build/env.sh go run build/ci.go install ./cmd/storm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/storm\" to launch storm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/storm.aar\" to use the library."

ios:
	build/env.sh go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Geth.framework\" to use the library."

test: all
	build/env.sh go run build/ci.go test

lint: ## Run linters.
	build/env.sh go run build/ci.go lint

clean:
	go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

storm-cross: storm-linux storm-darwin storm-windows storm-android storm-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/storm-*

storm-linux: storm-linux-386 storm-linux-amd64 storm-linux-arm storm-linux-mips64 storm-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/storm-linux-*

storm-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/storm
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/storm-linux-* | grep 386

storm-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/storm
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/storm-linux-* | grep amd64

storm-linux-arm: storm-linux-arm-5 storm-linux-arm-6 storm-linux-arm-7 storm-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/storm-linux-* | grep arm

storm-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/storm
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/storm-linux-* | grep arm-5

storm-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/storm
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/storm-linux-* | grep arm-6

storm-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/storm
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/storm-linux-* | grep arm-7

storm-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/storm
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/storm-linux-* | grep arm64

storm-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/storm
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/storm-linux-* | grep mips

storm-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/storm
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/storm-linux-* | grep mipsle

storm-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/storm
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/storm-linux-* | grep mips64

storm-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/storm
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/storm-linux-* | grep mips64le

storm-darwin: storm-darwin-386 storm-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/storm-darwin-*

storm-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/storm
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/storm-darwin-* | grep 386

storm-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/storm
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/storm-darwin-* | grep amd64

storm-windows: storm-windows-386 storm-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/storm-windows-*

storm-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/storm
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/storm-windows-* | grep 386

storm-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/storm
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/storm-windows-* | grep amd64
