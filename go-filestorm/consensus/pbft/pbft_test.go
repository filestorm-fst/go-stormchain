// Copyright 2019 The go-filestorm Authors
// This file is part of the go-filestorm library.
//
// The go-filestorm library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-filestorm library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-filestorm library. If not, see <http://www.gnu.org/licenses/>.

// Package pbft implements the practical Byzantine fault tolerance consensus engine.

package pbft

import (
	"os"

	"github/filestorm/go-filestorm/consensus/pbft/network"
)

func main() {
	nodeID := os.Args[1]
	server := network.NewServer(nodeID)

	server.Start()
}
