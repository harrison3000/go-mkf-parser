// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import (
	"fmt"
	"regexp"
	"regexp/syntax"
	"strconv"
	"strings"
)

// matcher typers
var (
	empty    = regexp.MustCompile(`^"()"`)
	regSpace = regexp.MustCompile(`^([\t ]+)`)

	literal    = regexp.MustCompile(`^"([^"\p{Cc}]+)"`)
	singleRune = regexp.MustCompile(`^'([^\p{C}])'`)
	simpleHex  = regexp.MustCompile(`^'([0-9A-F]{4,5})'`)
	tenHex     = regexp.MustCompile(`^'(10[0-9A-F]{4})'`)
	regdot     = regexp.MustCompile(`^(\.)`)
	regminus   = regexp.MustCompile(`^(-)`)
	regReg     = regexp.MustCompile(`^/((\\/|[^/])*)/`)

	regKnot = regexp.MustCompile(`^(§)`)

	regWhat  = regexp.MustCompile(`^(\?)`)
	regPlus  = regexp.MustCompile(`^(\+)`)
	regStar  = regexp.MustCompile(`^(\*)`)
	regRange = regexp.MustCompile(`^({\d+(,\d+)?})`)
)

type tokenKind rune

const (
	tkInvalid      tokenKind = 0
	tkEmpty        tokenKind = 'E'
	tkLiteral      tokenKind = 'L'
	tkSingleton    tokenKind = 'S'
	tkDot          tokenKind = '.'
	tkMinus        tokenKind = '-'
	tkRegex        tokenKind = 'r'
	tkRule         tokenKind = 'R'
	tkRuleRange    tokenKind = '#'
	tkRuleOperator tokenKind = '§'
	tkWhiteSpace   tokenKind = ' '
)

type altToken struct {
	val  string
	kind tokenKind
}

func str2alt(s string, allowEmpty bool) (alternative, error) {
	tks, e := tokenizeAlternative(s)
	if e != nil {
		return alternative{}, e
	}
	if allowEmpty && len(tks) == 1 && tks[0].kind == tkEmpty {
		return alternative{
			itens: []item{
				{kind: itemEmpty},
			},
		}, nil
	}

	var itens []item
	push := func(i item) {
		itens = append(itens, i)
	}

	for i := 0; i < len(tks); i++ {
		v := &tks[i]

		switch v.kind {
		case tkEmpty:
			return alternative{}, fmt.Errorf("unallowed empty found")
		case tkLiteral:
			//TODO possible optimization: single char strings -> singleton

			push(item{
				kind: itemLiteral,
				lit:  v.val,
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
			r, e := regexp.Compile(unescaped)
			if e != nil {
				return alternative{}, fmt.Errorf("error compiling regex: %w", e)
			}
			if !goodRegex(unescaped) {
				return alternative{}, fmt.Errorf("regexes must be anchored at the begining (^)")
			}

			push(item{
				kind:  itemRegex,
				regex: r,
			})

		case tkRule:
			push(item{
				kind: itemRule,
				lit:  v.val,
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

	consume := func(regex *regexp.Regexp, tKind tokenKind) bool {
		val, rest, ok := consumeRegex(s, regex)
		if !ok {
			return false
		}

		tks = append(tks, altToken{
			kind: tKind,
			val:  val,
		})
		s = rest
		return true
	}

	const maxTokens = 128
	for i := 0; i < maxTokens; i++ {
		switch {
		case
			consume(empty, tkEmpty),
			consume(regSpace, tkWhiteSpace),
			consume(literal, tkLiteral),
			consume(singleRune, tkSingleton),
			consume(simpleHex, tkSingleton),
			consume(tenHex, tkSingleton),
			consume(regdot, tkDot),
			consume(regminus, tkMinus),
			consume(regReg, tkRegex),
			consume(ruleName, tkRule),

			consume(regWhat, tkRuleRange),
			consume(regRange, tkRuleRange),
			consume(regPlus, tkRuleRange),
			consume(regStar, tkRuleRange),

			consume(regKnot, tkRuleOperator):

		default:
			col := len(orig) - len(s)
			return nil, fmt.Errorf("couldn't tokenize alternatives at column %d", col)
		}

		if isEmptyOrComment(s) {
			return validateAndFilterAltTokens(tks)
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
		ret.kind = itemRune
		r := tks[0].convertRune()
		ret.runes = runeRange{r, r}
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
		if len(tks) == 0 || tks[0].kind != tkMinus {
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
		kind:  itemSimpleRange,
		runes: base,
	}
	if len(excludes) != 0 {
		cplx := newComplexRange(base, excludes)
		if cplx == nil {
			return item{}, 0, fmt.Errorf("invalid exclusion range")
		}
		i = item{
			kind:         itemComplexRange,
			complexRange: cplx,
		}
	}

	return i, ol - len(tks), nil
}

func isSingleton(tks []altToken) bool {
	if len(tks) < 2 {
		return true
	}

	t := tks[1].kind
	return t != tkDot
}

func isRange(tks []altToken) bool {
	if len(tks) < 3 {
		return false
	}
	a := tks[0].kind == tkSingleton
	b := tks[1].kind == tkDot
	c := tks[2].kind == tkSingleton

	return a && b && c
}

func (a *alternative) isEmpty() bool {
	return len(a.itens) == 1 && a.itens[0].kind == itemEmpty
}

func goodRegex(s string) bool {
	//this only runs after the compilation, so we don't have to check errors
	rg, _ := syntax.Parse(s, syntax.Perl)
	si := rg.Simplify()
	prog, _ := syntax.Compile(si)

	sc := prog.StartCond()

	return sc == syntax.EmptyBeginText
}

// validateAndFilterAltTokens uses forbidden techniques to detect
// if a alternative is valid and removes the whitespaces
func validateAndFilterAltTokens(tks []altToken) ([]altToken, error) {
	var rr []rune
	for _, v := range tks {
		rr = append(rr, rune(v.kind))
	}
	synt := string(rr) + " "
	syntf := synt

	good := []string{
		" L ", " r ",
		" R§R ", " R§S ",
		" R# ", " R ",
		" S . S ", " S.S ", " S ",
		" - ",
		" E ", //this should be a special case, only one empty is allowed
	}

	for _, g := range good {
		rep := " ! "
		//we do replace twice because the spaces must overlap
		//we don't use stringsReplacer because the order matters
		//it must be each little thing twice, not the whole thing twice
		a := strings.ReplaceAll(syntf, g, rep)
		syntf = strings.ReplaceAll(a, g, rep)
	}

	if strings.ContainsRune(syntf, '§') {
		return nil, fmt.Errorf("misuse of the § operator")
	}

	if strings.Trim(syntf, "! ") != "" {
		return nil, fmt.Errorf("unrecognized alternative")
	}

	var ret []altToken
	for _, t := range tks {
		if t.kind != tkWhiteSpace {
			ret = append(ret, t)
		}
	}

	return ret, nil
}
