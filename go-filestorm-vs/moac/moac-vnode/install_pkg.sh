#!/bin/bash
# Install necessary packages for MOAC VNODE 
# Create fake Go workspace under build if it doesn't exist yet.
workspace="$PWD/build/_workspace"
# Set up the environment to use the workspace.
GOPATH="$workspace"
echo $GOPATH
export GOPATH
echo Set GOPATH $GOPATH
#install the following packages
go get github.com/aristanetworks/goarista/monotime
go get github.com/go-stack/stack
go get github.com/rcrowley/go-metrics
go get github.com/syndtr/goleveldb/leveldb
go get golang.org/x/net/context
go get google.golang.org/grpc
go get google.golang.org/grpc/reflection
go get gopkg.in/karalabe/cookiejar.v2/collections/prque
go get github.com/golang/protobuf/proto
go get github.com/hashicorp/golang-lru
go get github.com/pborman/uuid
go get github.com/rjeczalik/notify
go get github.com/patrickmn/go-cache
go get github.com/beevik/ntp
go get github.com/patrickmn/go-cache
go get gopkg.in/olebedev/go-duktape.v3
