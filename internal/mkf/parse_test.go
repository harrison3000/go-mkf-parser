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

func TestS2A(t *testing.T) {
	emp, e := str2alt(`"" //comment`, true)
	if e != nil {
		t.Fatal("error")
	}
	if len(emp) != 1 {
		t.Fatal("Wrong size")
	}
	if emp[0].typ != mtEmpty {
		t.Fatal("Wrong thing")
	}

	alt, e := str2alt(`"hello" 'a'`, false)

	if len(alt) != 2 {
		t.Fatal("Wrong size")
	}
	if e != nil {
		t.Fatal("error")
	}
	mtExp := matcher{typ: mtLiteral, lit: "hello"}
	if alt[0] != mtExp {
		t.Fatal("didn't parse literal")
	}

	mtExp = matcher{typ: mtRune, r: 'a'}
	if alt[1] != mtExp {
		t.Fatal("didn't parse rune")
	}

}
