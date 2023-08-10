// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// matcher typers
var (
	empty      = regexp.MustCompile(`^"()"`)
	literal    = regexp.MustCompile(`^"([^"\p{Cc}]+)"`)
	singleRune = regexp.MustCompile(`^'([^\p{C}])'`)
	simpleHex  = regexp.MustCompile(`^'([0-9A-F]{4,5})'`)
	tenHex     = regexp.MustCompile(`^'(10[0-9A-F]{4})'`)
	regdot     = regexp.MustCompile(`^(\.)`)
	regminus   = regexp.MustCompile(`^(-)`)
	regReg     = regexp.MustCompile(`^/((\\/|[^/])*)/`)
)

type tokensTy int

const (
	tkInvalid tokensTy = iota
	tkEmpty
	tkLiteral
	tkSingleton
	tkDot
	tkMinus
	tkRegex
	tkRule
)

type altToken struct {
	typ tokensTy
	val string
}

func str2alt(s string, allowEmpty bool) (alternative, error) {
	tks, e := tokenizeAlternative(s)
	if e != nil {
		return alternative{}, e
	}
	if allowEmpty && len(tks) == 1 && tks[0].typ == tkEmpty {
		return alternative{
			empty: true,
			itens: []item{
				{typ: itemEmpty},
			},
		}, nil
	}

	var itens []item
	push := func(i item) {
		itens = append(itens, i)
	}

	for i := 0; i < len(tks); i++ {
		v := &tks[i]

		switch v.typ {
		case tkEmpty:
			return alternative{}, fmt.Errorf("unallowed empty found")
		case tkLiteral:
			//TODO possible optimization: single char strings -> singleton

			push(item{
				typ: itemLiteral,
				lit: v.val,
			})
		case tkSingleton:
			it, skip, err := tksToItem(tks[i:])
			if err != nil {
				return alternative{}, fmt.Errorf("error interpreting range: %w", err)
			}
			i += skip - 1
			push(it)

		case tkRegex:
			unescaped := strings.ReplaceAll(v.val, `\/`, "/")
			//TODO warn on regexes that don't start with ^
			r, e := regexp.Compile(unescaped)
			if e != nil {
				return alternative{}, fmt.Errorf("error compiling regex: %w", e)
			}
			push(item{
				typ:     itemRegex,
				complex: r,
			})

		case tkRule:
			push(item{
				typ: itemRule,
				lit: v.val,
			})

		default:
			return alternative{}, fmt.Errorf("unexpected token")
		}
	}

	return alternative{itens: itens}, nil
}

func tokenizeAlternative(s string) ([]altToken, error) {
	orig := s
	var tks []altToken

	consume := func(regex *regexp.Regexp, typ tokensTy) bool {
		val, rest, ok := consumeRegex(s, regex)
		if !ok {
			return false
		}

		tks = append(tks, altToken{
			typ: typ,
			val: val,
		})
		s = rest
		return true
	}

	const maxTokens = 20
	for i := 0; i < maxTokens; i++ {
		sut := strings.TrimLeft(s, " \t")
		if sut == s { //didn't trim
			col := len(orig) - len(s)
			return nil, fmt.Errorf("required space not found at column %d", col)
		}
		s = sut

		switch {
		case
			consume(empty, tkEmpty),
			consume(literal, tkLiteral),
			consume(singleRune, tkSingleton),
			consume(simpleHex, tkSingleton),
			consume(tenHex, tkSingleton),
			consume(regdot, tkDot),
			consume(regminus, tkMinus),
			consume(regReg, tkRegex),
			consume(ruleName, tkRule):

		default:
			col := len(orig) - len(s)
			return nil, fmt.Errorf("couldn't tokenize alternatives at column %d", col)
		}

		if isEmptyOrComment(s) {
			return tks, nil
		}
	}

	return nil, fmt.Errorf("alternative too big (max: %d tokens)", maxTokens)
}

func (tk *altToken) convertRune() rune {
	s := tk.val

	if rn := []rune(s); len(rn) == 1 {
		return rn[0]
	}
	num, err := strconv.ParseInt(s, 16, 0)
	if err != nil {
		//the hexadecimal regexes should ensure this
		//part of the code isn't reached
		panic("the impossible became possible somehow")
	}

	return rune(num)
}

func tksToItem(tks []altToken) (item, int, error) {
	if isSingleton(tks) {
		var ret item
		ret.typ = itemRune
		ret.runes[0] = tks[0].convertRune()
		return ret, 1, nil
	}
	if !isRange(tks) {
		return item{}, 99999, fmt.Errorf("invalid syntax")
	}

	ol := len(tks)
	base := runeRange{
		tks[0].convertRune(),
		tks[2].convertRune(),
	}

	if !base.valid() {
		return item{}, 99999, fmt.Errorf("invalid range")
	}

	consume := func(i int) {
		tks = tks[i:]
	}

	consume(3) //the original range

	var excludes []runeRange
	var inception func() error

	inception = func() error {
		if len(tks) == 0 || tks[0].typ != tkMinus {
			return nil
		}
		consume(1) //the minus

		switch {
		case isSingleton(tks):
			r0 := tks[0].convertRune()
			excludes = append(excludes, runeRange{r0, r0})
			consume(1)
		case isRange(tks):
			r0, r1 := tks[0].convertRune(), tks[2].convertRune()
			excludes = append(excludes, runeRange{r0, r1})
			consume(3)
		default:
			return fmt.Errorf("invalid syntax, minus followed by wrong thing")
		}

		return inception() //we must go deeper
	}

	if err := inception(); err != nil {
		return item{}, 99999, err
	}

	i := item{
		typ:   itemSimpleRange,
		runes: base,
	}
	if len(excludes) != 0 {
		cplx := newComplexRange(base, excludes)
		if cplx == nil {
			return item{}, 0, fmt.Errorf("invalid exclusion range")
		}
		i = item{
			typ:     itemComplexRange,
			complex: cplx,
		}
	}

	return i, ol - len(tks), nil
}

func isSingleton(tks []altToken) bool {
	if len(tks) < 2 {
		return true
	}

	t := tks[1].typ
	return t != tkDot
}

func isRange(tks []altToken) bool {
	if len(tks) < 3 {
		return false
	}
	a := tks[0].typ == tkSingleton
	b := tks[1].typ == tkDot
	c := tks[2].typ == tkSingleton

	return a && b && c
}
