// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import "testing"

func TestBasic(t *testing.T) {
	p, e := NewParser(`
//just a test
test
    "wow"

tost
	'a'
`)
	if e != nil {
		t.Fatal("Error:", e)
	}
	if len(p.rules) == 0 {
		t.Fatal("shouldn't be empty")
	}

	_ = p
}

func TestS2ATypes(t *testing.T) {
	doTest := func(alt string, allowEmpty bool, types ...matcherType) {
		res, e := str2alt(alt, allowEmpty)
		if e != nil {
			t.Errorf("Error parsing alternative: %s", e)
			return
		}
		if len(res) != len(types) {
			t.Errorf("Wrong size")
		}
		for k, v := range res {
			if v.typ != types[k] {
				t.Error("Wrong type")
			}
		}
	}

	doTest(`"" //comment`, true, mtEmpty)

	doTest(`"hello" 'a' 'f' . 't'`, false,
		mtLiteral,
		mtRune,
		mtSimpleRange,
	)

	_, e := str2alt(`"hello" ""`, false)
	if e == nil {
		t.Error("should be error")
	}
}
