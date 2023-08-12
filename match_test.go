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
	digit value //TODO tail call optimization (?)
	digit

digit
	'0' . '9'

ws
	""
	/^\s+/
	`)

	if e != nil {
		t.Fatalf("Should be nil: %s", e)
	}

	mustGoRight := func(s string) *Node {
		res, e := p.ParseString(s)
		t.Logf("Now testing: %s", s)
		if e != nil {
			t.Error("Shouldn't have failed: ", e)
		} else if res == nil {
			t.Error("Shouldn't be nil result")
		} else if res.val != s {
			t.Errorf("Wrong node value, expected: %v, got %v", s, res.val)
		} else {
			t.Log("Ok")
		}
		return res
	}

	mustGoRight("[7]")
	mustGoRight("[729]")
	mustGoRight("[72,899]")
	mustGoRight("[720,444,22,123,5, 123 ,123]")

	_, e = p.ParseString("[ 720,444,22,123,5, 1z23 ,12]")
	if e == nil {
		t.Error("Should have failed")
	}

}

//TODO complex range test
