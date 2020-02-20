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
	"testing"
)

func TestURLParsing(t *testing.T) {
	url, err := parseURL("https://filestorm.org")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if url.Scheme != "https" {
		t.Errorf("expected: %v, got: %v", "https", url.Scheme)
	}
	if url.Path != "filestorm.org" {
		t.Errorf("expected: %v, got: %v", "filestorm.org", url.Path)
	}

	_, err = parseURL("filestorm.org")
	if err == nil {
		t.Error("expected err, got: nil")
	}
}

func TestURLString(t *testing.T) {
	url := URL{Scheme: "https", Path: "filestorm.org"}
	if url.String() != "https://filestorm.org" {
		t.Errorf("expected: %v, got: %v", "https://filestorm.org", url.String())
	}

	url = URL{Scheme: "", Path: "filestorm.org"}
	if url.String() != "filestorm.org" {
		t.Errorf("expected: %v, got: %v", "filestorm.org", url.String())
	}
}

func TestURLMarshalJSON(t *testing.T) {
	url := URL{Scheme: "https", Path: "filestorm.org"}
	json, err := url.MarshalJSON()
	if err != nil {
		t.Errorf("unexpcted error: %v", err)
	}
	if string(json) != "\"https://filestorm.org\"" {
		t.Errorf("expected: %v, got: %v", "\"https://filestorm.org\"", string(json))
	}
}

func TestURLUnmarshalJSON(t *testing.T) {
	url := &URL{}
	err := url.UnmarshalJSON([]byte("\"https://filestorm.org\""))
	if err != nil {
		t.Errorf("unexpcted error: %v", err)
	}
	if url.Scheme != "https" {
		t.Errorf("expected: %v, got: %v", "https", url.Scheme)
	}
	if url.Path != "filestorm.org" {
		t.Errorf("expected: %v, got: %v", "https", url.Path)
	}
}

func TestURLComparison(t *testing.T) {
	tests := []struct {
		urlA   URL
		urlB   URL
		expect int
	}{
		{URL{"https", "filestorm.org"}, URL{"https", "filestorm.org"}, 0},
		{URL{"http", "filestorm.org"}, URL{"https", "filestorm.org"}, -1},
		{URL{"https", "filestorm.org/a"}, URL{"https", "filestorm.org"}, 1},
		{URL{"https", "abc.org"}, URL{"https", "filestorm.org"}, -1},
	}

	for i, tt := range tests {
		result := tt.urlA.Cmp(tt.urlB)
		if result != tt.expect {
			t.Errorf("test %d: cmp mismatch: expected: %d, got: %d", i, tt.expect, result)
		}
	}
}
