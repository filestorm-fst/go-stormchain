package types

import (
	pb "github.com/filestorm/go-filestorm/moac/chain3go/lib/proto"
)

type ScsServerConnection struct {
	ScsHostAddress string
	ScsId          string
	LiveFlag       bool
	Stream         *pb.Vnode_ScsPushServer
	Req            chan *pb.ScsPushMsg
	Cancel         chan bool
	RetryCount     uint
}
