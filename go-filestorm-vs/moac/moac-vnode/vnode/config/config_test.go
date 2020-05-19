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

package config

import (
	"fmt"
	"reflect"
	"testing"
)

// Tests that Configuration files
func TestCheckGetConfiguration(t *testing.T) {
	type test struct {
		readconfig *Configuration
		infilepath string
		wantErr    error
		SCSService bool
	}
	tests := []test{
		{readconfig: &Configuration{}, infilepath: "./error.json", wantErr: nil, SCSService: false},
		{readconfig: nil, infilepath: "./vnodeconfig.json", wantErr: nil, SCSService: true},
	}

	for _, test := range tests {
		readConfig, err := GetConfiguration(test.infilepath)
		fmt.Printf("Read config %v\n", readConfig.SCSService)
		if !reflect.DeepEqual(err, test.wantErr) {
			t.Errorf("error mismatch:\nRead input: %v\nGet err: %v\nwant: %v", test.infilepath, err, test.wantErr)
		}
		if !reflect.DeepEqual(readConfig.SCSService, test.SCSService) {
			t.Errorf("error mismatch:\nRead input: %v\nGet value: %v\nwant: %v", test.infilepath, readConfig.SCSService, test.SCSService)
		}
	}
}
