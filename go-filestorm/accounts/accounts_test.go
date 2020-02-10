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

package accounts

import (
	"bytes"
	"testing"

	"github.com/filestorm/go-filestorm/common/hexutil"
)

func TestTextHash(t *testing.T) {
	hash := TextHash([]byte("Hello FileStorm"))
	want := hexutil.MustDecode("0xe82707262ca0aacbd8daf253cebcf05ad6b898223b2a445460b5a0d3fe5d7912")
	if !bytes.Equal(hash, want) {
		t.Fatalf("wrong hash: %x", hash)
	}
}
