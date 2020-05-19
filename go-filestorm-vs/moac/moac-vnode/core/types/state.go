// Copyright 2017 The moac-vnode Authors
// This file is part of the moac-vnode library.
//
// The go-moac-vnode library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-moac-vnode library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the moac-vnode library. If not, see <http://www.gnu.org/licenses/>.

package types

type AccountChgDump struct {
	Root     string                      `json:"root"`
	Accounts map[string]AccountChgStruct `json:"accounts"`
}

type AccountChgStruct struct {
	Address    string            `json:"address"`
	BalanceChg string            `json:"balanceChg"`
	Nonce      uint64            `json:"nonce"`
	Root       string            `json:"root"`
	CodeHash   string            `json:"codeHash"`
	Code       string            `json:"code"`
	Storage    map[string]string `json:"storage"`
}
