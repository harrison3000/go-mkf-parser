// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import "testing"

const testArrayParser = `
rootRule
	array

array
	'[' values ']'

values
	arrElement
	arrElement ',' values

arrElement
	ws decValue ws
	ws hexValue ws
	ws array	ws

hexValue
	/^0x[A-Fa-f0-9]+/

decValue
	digit decValue
	digit

digit
	'0' . '9'

ws
	""
	/^\s+/
	`

func mustGoAlright(p *Parser, t *testing.T, s string) *Node {
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

func TestMatch(t *testing.T) {
	p, e := NewParser(testArrayParser)

	if e != nil {
		t.Fatalf("Should be nil: %s", e)
	}

	mustGoRight := func(s string) *Node {
		res := mustGoAlright(p, t, s)
		return res
	}

	mustFail := func(s string) *Node {
		res, e := p.ParseString(s)
		if e == nil {
			t.Error("Should have failed")
		}
		return res
	}

	mustGoRight("[7]")
	mustGoRight("[729]")
	mustGoRight("[72,899]")
	mustGoRight("[720,444,22,123,5, 123 ,123]")

	mustFail("[ 720,444,22,123,5, 1z23 ,12]")

	p, e = NewParser(`
rootRule
	'a' . 'z' - 'p' - 'd' . 'f' - l
	"literal"
	`)

	if e != nil {
		t.Fatalf("Should be nil: %s", e)
	}

	mustGoRight("a")
	mustGoRight("t")
	mustGoRight("q")
	mustGoRight("z")
	mustGoRight("literal")

	mustFail("A")
	mustFail("p")
	mustFail("f")
	mustFail("e")
	mustFail("d")
	mustFail("l")

}

func BenchmarkParsing(b *testing.B) {
	p, e := NewParser(testArrayParser)
	if e != nil {
		b.Fatalf("Error compiling grammar: %s", e)
	}

	mustGoRight := func(s string) {
		r, e := p.ParseString(s)
		if r == nil || e != nil {
			b.Fatal("failed parsing ")
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mustGoRight(`[720,444,22,123,5, 123 ,123]`)
		mustGoRight(`[720, 4654844, 287982, 123, 5546, 0xccc123 ,0xfba123]`)
		mustGoRight(`[720, 4654844, 287982, 
		[1234125,0x1234125ab61234,12341243123,0x123f51265134652,2345234562345,234623452345],
		123, 5546, 0xccc123 ,0xfba123]`)
	}
}

func TestComplexes(t *testing.T) {
	p, e := NewParser(`
csv
	digits§','
	digits§commaWs

commaWs
	/^\s*,\s*/

digits
	digit+

digit
	'0' . '9'
`)

	mustFail := func(s string) *Node {
		res, e := p.ParseString(s)
		if e == nil {
			t.Error("Should have failed")
		}
		return res
	}

	if e != nil {
		t.Fatalf("Failed creating parser, should be nil: %s", e)
	}

	mustGoAlright(p, t, "123123,123123,412343445, 346346,34563456 , 3456")
	mustGoAlright(p, t, "123123,123123,41234")
	mustGoAlright(p, t, "123123 , 123123 , 41234")

	mustFail("123123,123123,")

}
