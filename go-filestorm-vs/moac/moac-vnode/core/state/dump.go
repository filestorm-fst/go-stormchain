// Copyright 2014 The MOAC-core Authors
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

package state

import (
	"encoding/json"
	"fmt"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/rlp"
	"github.com/filestorm/go-filestorm/moac/moac-lib/trie"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/vnode/types"
	pb "github.com/filestorm/go-filestorm/moac/moac-lib/proto"
	libtypes "github.com/filestorm/go-filestorm/moac/moac-lib/types"
)

type DumpAccount struct {
	Balance  string            `json:"balance"`
	Nonce    uint64            `json:"nonce"`
	Root     string            `json:"root"`
	CodeHash string            `json:"codeHash"`
	Code     string            `json:"code"`
	Storage  map[string]string `json:"storage"`
}

type Dump struct {
	Root     string                 `json:"root"`
	Accounts map[string]DumpAccount `json:"accounts"`
}

func (self *StateDB) RawDump() Dump {
	dump := Dump{
		Root:     fmt.Sprintf("%x", self.trie.Hash()),
		Accounts: make(map[string]DumpAccount),
	}

	it := trie.NewIterator(self.trie.NodeIterator(nil))
	for it.Next() {
		addr := self.trie.GetKey(it.Key)
		var data Account
		if err := rlp.DecodeBytes(it.Value, &data); err != nil {
			panic(err)
		}

		obj := newObject(nil, common.BytesToAddress(addr), data, nil)
		account := DumpAccount{
			Balance:  data.Balance.String(),
			Nonce:    data.Nonce,
			Root:     common.Bytes2Hex(data.Root[:]),
			CodeHash: common.Bytes2Hex(data.CodeHash),
			Code:     common.Bytes2Hex(obj.Code(self.db)),
			Storage:  make(map[string]string),
		}
		storageIt := trie.NewIterator(obj.getTrie(self.db).NodeIterator(nil))
		for storageIt.Next() {
			account.Storage[common.Bytes2Hex(self.trie.GetKey(storageIt.Key))] = common.Bytes2Hex(storageIt.Value)
		}
		dump.Accounts[common.Bytes2Hex(addr)] = account
	}
	return dump
}

func (self *StateDB) Dump() []byte {
	json, err := json.MarshalIndent(self.RawDump(), "", "    ")
	if err != nil {
		fmt.Println("dump err", err)
	}

	return json
}

func (self *StateDB) DumpAccountStorage(addrin common.Address) []byte {

	it := trie.NewIterator(self.trie.NodeIterator(nil))
	for it.Next() {
		addr := self.trie.GetKey(it.Key)
		if addrin == common.BytesToAddress(addr) {
			var data Account
			if err := rlp.DecodeBytes(it.Value, &data); err != nil {
				panic(err)
			}

			obj := newObject(nil, common.BytesToAddress(addr), data, nil)
			account := DumpAccount{
				Balance:  data.Balance.String(),
				Nonce:    data.Nonce,
				Root:     common.Bytes2Hex(data.Root[:]),
				CodeHash: common.Bytes2Hex(data.CodeHash),
				Code:     "", //common.Bytes2Hex(obj.Code(self.db)),
				Storage:  make(map[string]string),
			}
			storageIt := trie.NewIterator(obj.getTrie(self.db).NodeIterator(nil))
			for storageIt.Next() {
				account.Storage[common.Bytes2Hex(self.trie.GetKey(storageIt.Key))] = common.Bytes2Hex(storageIt.Value)
			}

			json, _ := json.MarshalIndent(account, "", "    ")
			return json
		}

	}

	return nil
}

func (self *StateDB) DumpContractStorage(addrin common.Address, request []*pb.StorageRequest) []byte {

	it := trie.NewIterator(self.trie.NodeIterator(nil))
	for it.Next() {
		addr := self.trie.GetKey(it.Key)
		if addrin == common.BytesToAddress(addr) {
			var data Account
			if err := rlp.DecodeBytes(it.Value, &data); err != nil {
				panic(err)
			}

			obj := newObject(nil, common.BytesToAddress(addr), data, nil)
			account := types.ContractInfo{
				Balance:  data.Balance,
				Nonce:    data.Nonce,
				Root:     data.Root,
				CodeHash: data.CodeHash,
				Code:     obj.Code(self.db),
				Storage:  make(map[string]string),
			}
			
			storage := make(map[string]string)
			storageIt := trie.NewIterator(obj.getTrie(self.db).NodeIterator(nil))
			for storageIt.Next() {
				storage[common.Bytes2Hex(self.trie.GetKey(storageIt.Key))] = common.Bytes2Hex(storageIt.Value)
				//log.Info("key:val", common.Bytes2Hex(self.trie.GetKey(storageIt.Key)),  common.Bytes2Hex(storageIt.Value))
			}

			// for _, val := range request {
			// 	key := common.Bytes2Hex(val.Storagekey)
			// 	position := common.Bytes2Hex(val.Position)
			// 	structformat := val.Structformat
			// 	switch val.Reqtype {
			// 	case 0:
			// 		for k, value := range storage {
			// 			account.Storage[k] = value
			// 		}
			// 	case 1:
			// 		if len(position) == 0 {
			// 			var num int64
			// 			strlen := storage[key]
			// 			if len(strlen) > 2 {
			// 				num, _ = strconv.ParseInt(strlen[2:], 16, 64)
			// 			} else {
			// 				num, _ = strconv.ParseInt(strlen, 16, 64)
			// 			}
			// 			account.Storage[key] = storage[key]
			// 			keys := common.KeytoKey(key)
			// 			for i := int64(0); i < num; i++ {
			// 				if len(structformat) != 0 {
			// 					key0 := keys
			// 					for j := 0; j < len(structformat); j++ {
			// 						if structformat[j] == '1' {
			// 							account.Storage[key0] = storage[key0]
			// 						} else if structformat[j] == '2' {
			// 							account.Storage[key0] = storage[key0]
			// 							var num0 int64
			// 							strlen0 := storage[key0]
			// 							if len(strlen0) > 2 {
			// 								num0, _ = strconv.ParseInt(strlen0[2:], 16, 64)
			// 							} else {
			// 								num0, _ = strconv.ParseInt(strlen0, 16, 64)
			// 							}
			// 							key1 := common.KeytoKey(key0)
			// 							for k := int64(0); k < num0; k++ {
			// 								account.Storage[key1] = storage[key1]
			// 								key1 = common.IncreaseHexByOne(key1)
			// 							}
			// 						} else if structformat[j] == '3' {
			// 							nlen := len(storage[key0])
			// 							if nlen == 66 {
			// 								account.Storage[key0] = storage[key0]
			// 							} else if nlen == 2 {
			// 								account.Storage[key0] = storage[key0]
			// 								key1 := common.KeytoKey(key0)
			// 								account.Storage[key1] = storage[key1]
			// 								key1 = common.IncreaseHexByOne(key1)
			// 								account.Storage[key1] = storage[key1]
			// 							} else if nlen > 2 && nlen < 66 {
			// 								account.Storage[key0] = storage[key0]
			// 								num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
			// 								key1 := common.KeytoKey(key0)
			// 								for i := num -1; i > 0; {
			// 									account.Storage[key1] = storage[key1]
			// 									key1 = common.IncreaseHexByOne(key1)
			// 									i = i - 64
			// 								}
			// 							}
			// 						}
			// 						key0 = common.IncreaseHexByOne(key0)
			// 					}
			// 				} else {
			// 					// account.Storage[keys] = storage[keys]
			// 					key0 := keys
			// 					nlen := len(storage[key0])
			// 					if nlen == 66 {
			// 						account.Storage[key0] = storage[key0]
			// 					} else if nlen == 2 {
			// 						account.Storage[key0] = storage[key0]
			// 						key1 := common.KeytoKey(key0)
			// 						account.Storage[key1] = storage[key1]
			// 						key1 = common.IncreaseHexByOne(key1)
			// 						account.Storage[key1] = storage[key1]
			// 					} else if nlen > 2 && nlen < 66 {
			// 						account.Storage[key0] = storage[key0]
			// 						if nlen < 7 {
			// 							num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
			// 							key1 := common.KeytoKey(key0)
			// 							for i := num -1; i > 0; {
			// 								if storage[key1] != "" {
			// 									account.Storage[key1] = storage[key1]
			// 								}
			// 								key1 = common.IncreaseHexByOne(key1)
			// 								i = i - 64
			// 							}
			// 						}
			// 					}
			// 				}
			// 				keys = common.IncreaseHexByOne(keys)
			// 			}
	
			// 		} else {
			// 			num, _ := strconv.ParseInt(position, 16, 64)
			// 			keys := common.KeytoKey(key)
			// 			keys = common.IncreaseHexByNum(num, keys)
			// 			if len(structformat) != 0 {
			// 				key0 := keys
			// 				for j := 0; j < len(structformat); j++ {
			// 					if structformat[j] == '1' {
			// 						account.Storage[key0] = storage[key0]
			// 					} else if structformat[j] == '2' {
			// 						account.Storage[key0] = storage[key0]
			// 						var num0 int64
			// 						strlen0 := storage[key0]
			// 						if len(strlen0) > 2 {
			// 							num0, _ = strconv.ParseInt(strlen0[2:], 16, 64)
			// 						} else {
			// 							num0, _ = strconv.ParseInt(strlen0, 16, 64)
			// 						}
			// 						key1 := common.KeytoKey(key0)
			// 						for k := int64(0); k < num0; k++ {
			// 							account.Storage[key1] = storage[key1]
			// 							key1 = common.IncreaseHexByOne(key1)
			// 						}
			// 					} else if structformat[j] == '3' {
			// 						nlen := len(storage[key0])
			// 						if nlen == 66 {
			// 							account.Storage[key0] = storage[key0]
			// 						} else if nlen == 2 {
			// 							account.Storage[key0] = storage[key0]
			// 							key1 := common.KeytoKey(key0)
			// 							account.Storage[key1] = storage[key1]
			// 							key1 = common.IncreaseHexByOne(key1)
			// 							account.Storage[key1] = storage[key1]
			// 						} else if nlen > 2 && nlen < 66 {
			// 							account.Storage[key0] = storage[key0]
			// 							num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
			// 							key1 := common.KeytoKey(key0)
			// 							for i := num -1; i > 0; {
			// 								account.Storage[key1] = storage[key1]
			// 								key1 = common.IncreaseHexByOne(key1)
			// 								i = i - 64
			// 							}
			// 						}
			// 					}
			// 					key0 = common.IncreaseHexByOne(key0)
			// 				}
			// 			} else {
			// 				// account.Storage[keys] = storage[keys]
			// 				key0 := keys
			// 				nlen := len(storage[key0])
			// 				if nlen == 66 {
			// 					account.Storage[key0] = storage[key0]
			// 				} else if nlen == 2 {
			// 					account.Storage[key0] = storage[key0]
			// 					key1 := common.KeytoKey(key0)
			// 					account.Storage[key1] = storage[key1]
			// 					key1 = common.IncreaseHexByOne(key1)
			// 					account.Storage[key1] = storage[key1]
			// 				} else if nlen > 2 && nlen < 66 {
			// 					account.Storage[key0] = storage[key0]
			// 					if nlen < 7 {
			// 						num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
			// 						key1 := common.KeytoKey(key0)
			// 						for i := num -1; i > 0; {
			// 							if storage[key1] != "" {
			// 								account.Storage[key1] = storage[key1]
			// 							}
			// 							key1 = common.IncreaseHexByOne(key1)
			// 							i = i - 64
			// 						}
			// 					}
			// 				}
			// 			}
			// 		}
			// 	case 2:
			// 		keys := common.KeytoKey(position + key)
			// 		if len(structformat) != 0 {
			// 			key0 := keys
			// 			for j := 0; j < len(structformat); j++ {
			// 				if structformat[j] == '1' {
			// 					account.Storage[key0] = storage[key0]
			// 				} else if structformat[j] == '2' {
			// 					account.Storage[key0] = storage[key0]
			// 					var num0 int64
			// 					strlen0 := storage[key0]
			// 					if len(strlen0) > 2 {
			// 						num0, _ = strconv.ParseInt(strlen0[2:], 16, 64)
			// 					} else {
			// 						num0, _ = strconv.ParseInt(strlen0, 16, 64)
			// 					}
	
			// 					key1 := common.KeytoKey(key0)
			// 					for k := int64(0); k < num0; k++ {
			// 						account.Storage[key1] = storage[key1]
			// 						key1 = common.IncreaseHexByOne(key1)
			// 					}
			// 				} else if structformat[j] == '3' {
			// 					nlen := len(storage[key0])
			// 					if nlen == 66 {
			// 						account.Storage[key0] = storage[key0]
			// 					} else if nlen == 2 {
			// 						account.Storage[key0] = storage[key0]
			// 						key1 := common.KeytoKey(key0)
			// 						account.Storage[key1] = storage[key1]
			// 						key1 = common.IncreaseHexByOne(key1)
			// 						account.Storage[key1] = storage[key1]
			// 					} else if nlen > 2 && nlen < 66 {
			// 						account.Storage[key0] = storage[key0]
			// 						num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
			// 						key1 := common.KeytoKey(key0)
			// 						for i := num -1; i > 0; {
			// 							account.Storage[key1] = storage[key1]
			// 							key1 = common.IncreaseHexByOne(key1)
			// 							i = i - 64
			// 						}
			// 					}
			// 				}
			// 				key0 = common.IncreaseHexByOne(key0)
			// 			}
			// 		} else {
			// 			// account.Storage[keys] = storage[keys]
			// 			key0 := keys
			// 			nlen := len(storage[key0])
			// 			if nlen == 66 {
			// 				account.Storage[key0] = storage[key0]
			// 			} else if nlen == 2 {
			// 				account.Storage[key0] = storage[key0]
			// 				key1 := common.KeytoKey(key0)
			// 				account.Storage[key1] = storage[key1]
			// 				key1 = common.IncreaseHexByOne(key1)
			// 				account.Storage[key1] = storage[key1]
			// 			} else if nlen > 2 && nlen < 66 {
			// 				account.Storage[key0] = storage[key0]
			// 				if nlen < 7 {
			// 					num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
			// 					key1 := common.KeytoKey(key0)
			// 					for i := num -1; i > 0; {
			// 						if storage[key1] != "" {
			// 							account.Storage[key1] = storage[key1]
			// 						}
			// 						key1 = common.IncreaseHexByOne(key1)
			// 						i = i - 64
			// 					}
			// 				}
			// 			}
			// 		}
			// 	case 3:
			// 		key0 := key
			// 		for j := 0; j < len(structformat); j++ {
			// 			if structformat[j] == '1' {
			// 				account.Storage[key0] = storage[key0]
			// 			} else if structformat[j] == '2' {
			// 				account.Storage[key0] = storage[key0]
			// 				var num0 int64
			// 				strlen0 := storage[key0]
			// 				if len(strlen0) > 2 {
			// 					num0, _ = strconv.ParseInt(strlen0[2:], 16, 64)
			// 				} else {
			// 					num0, _ = strconv.ParseInt(strlen0, 16, 64)
			// 				}
			// 				key1 := common.KeytoKey(key0)
			// 				for k := int64(0); k < num0; k++ {
			// 					account.Storage[key1] = storage[key1]
			// 					key1 = common.IncreaseHexByOne(key1)
			// 				}
			// 			} else if structformat[j] == '3' {
			// 				nlen := len(storage[key0])
			// 				if nlen == 66 {
			// 					account.Storage[key0] = storage[key0]
			// 				} else if nlen == 2 {
			// 					account.Storage[key0] = storage[key0]
			// 					key1 := common.KeytoKey(key0)
			// 					account.Storage[key1] = storage[key1]
			// 					key1 = common.IncreaseHexByOne(key1)
			// 					account.Storage[key1] = storage[key1]
			// 				} else if nlen > 2 && nlen < 66 {
			// 					account.Storage[key0] = storage[key0]
			// 					num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
			// 					key1 := common.KeytoKey(key0)
			// 					for i := num -1; i > 0; {
			// 						account.Storage[key1] = storage[key1]
			// 						key1 = common.IncreaseHexByOne(key1)
			// 						i = i - 64
			// 					}
			// 				}
			// 			}
			// 			key0 = common.IncreaseHexByOne(key0)
			// 		}
			// 	case 4:
			// 		account.Storage[key] = storage[key]
			// 	case 5:
			// 		nlen := len(storage[key])
			// 		if nlen == 66 {
			// 			account.Storage[key] = storage[key]
			// 		} else if nlen == 2 {
			// 			account.Storage[key] = storage[key]
			// 			key0 := common.KeytoKey(key)
			// 			account.Storage[key0] = storage[key0]
			// 			key0 = common.IncreaseHexByOne(key0)
			// 			account.Storage[key0] = storage[key0]
			// 		} else if nlen > 2 && nlen < 66 {
			// 			account.Storage[key] = storage[key]
			// 			num, _ := strconv.ParseInt(storage[key][2:], 16, 64)
			// 			key0 := common.KeytoKey(key)
			// 			for i := num -1; i > 0; {
			// 				account.Storage[key0] = storage[key0]
			// 				key0 = common.IncreaseHexByOne(key0)
			// 				i = i - 64
			// 			}
			// 		}
			// 	default:
	
			// 	}
			// }

			// for _, val := range request {
			// 	account.Storage = libtypes.ScreeningStorage(storage, val)
			// }

			account.Storage = libtypes.ScreeningStorage(storage, request)
			json, _ := json.Marshal(account)
			return json
		}

	}

	return nil
}
