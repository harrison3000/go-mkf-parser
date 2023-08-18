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

func TestS2AKinds(t *testing.T) {
	doTest := func(alt string, allowEmpty bool, kinds ...itemKind) {
		a, e := str2alt(" "+alt, allowEmpty)
		if e != nil {
			t.Errorf("Error parsing alternative: %s", e)
			return
		}
		if lg, le := len(a.itens), len(kinds); lg != le {
			t.Errorf("Wrong size, expected: %d, got: %d", le, lg)
			return
		}
		for k, v := range a.itens {
			if v.kind != kinds[k] {
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

	doTest(`"hi" 'b' . 'd' 'a' . 'z' - 't' . 'v' - 'h' /^a \/ aa/ '10ABCD'`, false,
		itemLiteral,
		itemSimpleRange,
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

func TestAltTokenizer(t *testing.T) {
	doTest := func(alt string, kinds ...tokenKind) {
		a, e := tokenizeAlternative(" " + alt)
		if e != nil {
			t.Errorf("Error parsing alternative: %s", e)
			return
		}
		if lg, le := len(a), len(kinds); lg != le {
			t.Errorf("Wrong size, expected: %d, got: %d", le, lg)
			return
		}
		for k, v := range a {
			if v.kind != kinds[k] {
				t.Error("Wrong type")
			}
		}
	}

	doTest(`aaa* aaa+ aaa? aaa{2,5} bbbb§iiii`,
		tkRule, tkRuleRange,
		tkRule, tkRuleRange,
		tkRule, tkRuleRange,
		tkRule, tkRuleRange,
		tkRule, tkRuleOperator, tkRule,
	)
}

func TestNotExists(t *testing.T) {
	_, e := NewParser(`
test
    anotherRule
`)
	if e == nil {
		t.Fatal("False positive")
	}
}

func BenchmarkCompilation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewParser(testArrayParser)
	}
}
