// requestData project syncing.go
package requestData

import "github.com/filestorm/go-filestorm/moac/chain3go/types"

type SyncingResponse struct {
	StartingBlock types.ComplexIntResponse `json:"startingBlock"`
	CurrentBlock  types.ComplexIntResponse `json:"currentBlock"`
	HighestBlock  types.ComplexIntResponse `json:"highestBlock"`
}
