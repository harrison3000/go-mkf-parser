// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import "testing"

func TestMatch(t *testing.T) {
	p, e := NewParser(`
rootRule
	array

array
	'[' values ']'

values
	arrElement
	arrElement ',' values

arrElement
	ws value ws

value
	'0' . '9'

ws
	""
	/^\s+/
	`)

	if e != nil {
		t.Fatalf("Should be nil: %s", e)
	}

	res, e := p.ParseString("[720,444,22,123,5, 123 ,123]")
	if e != nil {
		t.Error("Shouldn't have failed: ", e)
	}

	_ = res

	_, e = p.ParseString("[ 720,444,22,123,5, 1z23 ,12]")
	if e == nil {
		t.Error("Should have failed")
	}

}
