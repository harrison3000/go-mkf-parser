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
	doTest := func(alt string, allowEmpty bool, types ...itemType) {
		a, e := str2alt(" "+alt, allowEmpty)
		if e != nil {
			t.Errorf("Error parsing alternative: %s", e)
			return
		}
		if lg, le := len(a.itens), len(types); lg != le {
			t.Errorf("Wrong size, expected: %d, got: %d", le, lg)
		}
		for k, v := range a.itens {
			if v.typ != types[k] {
				t.Error("Wrong type")
			}
		}
	}

	doTest(`"" //comment`, true, itemEmpty)

	doTest(`"hello" 'a' 'f' . 't'`, false,
		itemLiteral,
		itemRune,
		itemSimpleRange,
	)

	doTest(`"hi" 'b' . 'd' 'a' . 'z' - 't' . 'v' - 'h' /a \/ aa/ '10ABCD'`, false,
		itemLiteral,
		itemComplexRange,
		itemRegex,
		itemRune,
	)

	doTestError := func(alt string, allowEmpty bool) {
		_, e := str2alt(alt, allowEmpty)
		if e == nil {
			t.Errorf("expected error")
		}
	}

	doTestError(` "hello" ""`, false) //empty not alone
	doTestError(` 'a'.'z'`, true)     //no space sparating
}
