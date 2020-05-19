// Copyright 2017 The MOAC-core Authors
// This file is part of the MOAC-core library.
//
// The MOAC-core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The MOAC-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the MOAC-core library. If not, see <http://www.gnu.org/licenses/>.

package tests

import (
	"fmt"
	"math/big"

	"github.com/filestorm/go-filestorm/moac/moac-vnode/params"
)

// This table defines supported forks and their chain config.
var Forks = map[string]*params.ChainConfig{
	// "Frontier": &params.ChainConfig{
	// 	ChainId: big.NewInt(1),
	// },
	"Pangu": &params.ChainConfig{
		ChainId:            big.NewInt(1),
		PanguBlock:         big.NewInt(0),
		RemoveEmptyAccount: true,
	},
	// "EIP150": &params.ChainConfig{
	// 	ChainId:     big.NewInt(1),
	// 	PanguBlock:  big.NewInt(0),
	// 	EIP150Block: big.NewInt(0),
	// },
	// "EIP158": &params.ChainConfig{
	// 	ChainId:     big.NewInt(1),
	// 	PanguBlock:  big.NewInt(0),
	// 	EIP150Block: big.NewInt(0),
	// 	EIP155Block: big.NewInt(0),
	// 	EIP158Block: big.NewInt(0),
	// },
	// "Byzantium": &params.ChainConfig{
	// 	ChainId:        big.NewInt(1),
	// 	PanguBlock:     big.NewInt(0),
	// 	EIP150Block:    big.NewInt(0),
	// 	EIP155Block:    big.NewInt(0),
	// 	EIP158Block:    big.NewInt(0),
	// 	DAOForkBlock:   big.NewInt(0),
	// 	ByzantiumBlock: big.NewInt(0),
	// },
	// "FrontierToPanguAt5": &params.ChainConfig{
	// 	ChainId:    big.NewInt(1),
	// 	PanguBlock: big.NewInt(5),
	// },
	// "PanguToEIP150At5": &params.ChainConfig{
	// 	ChainId:     big.NewInt(1),
	// 	PanguBlock:  big.NewInt(0),
	// 	EIP150Block: big.NewInt(5),
	// },
	// "PanguToDaoAt5": &params.ChainConfig{
	// 	ChainId:        big.NewInt(1),
	// 	PanguBlock:     big.NewInt(0),
	// 	DAOForkBlock:   big.NewInt(5),
	// 	DAOForkSupport: true,
	// },
	// "EIP158ToByzantiumAt5": &params.ChainConfig{
	// 	ChainId:        big.NewInt(1),
	// 	PanguBlock:     big.NewInt(0),
	// 	EIP150Block:    big.NewInt(0),
	// 	EIP155Block:    big.NewInt(0),
	// 	EIP158Block:    big.NewInt(0),
	// 	ByzantiumBlock: big.NewInt(5),
	// },
}

// UnsupportedForkError is returned when a test requests a fork that isn't implemented.
type UnsupportedForkError struct {
	Name string
}

func (e UnsupportedForkError) Error() string {
	return fmt.Sprintf("unsupported fork %q", e.Name)
}
