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
	literal    = regexp.MustCompile(`^"([^"]+)"`)
	singleRune = regexp.MustCompile(`^'(.)'`)
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
)

type altToken struct {
	typ tokensTy
	val string
	r   rune
}

func str2alt(s string, allowEmpty bool) (alternative, error) {
	tks, e := tokenizeAlternative(s)
	if e != nil {
		return nil, e
	}
	if allowEmpty && len(tks) == 1 && tks[0].typ == tkEmpty {
		return alternative{
			matcher{typ: mtEmpty},
		}, nil
	}

	for k := range tks {
		v := &tks[k]
		if v.typ == tkEmpty {
			return nil, fmt.Errorf("unallowed empty found")
		}
		v.convertRune()
	}

	_ = tks
	//TODO implement
	return nil, nil
}

func tokenizeAlternative(s string) ([]altToken, error) {
	var tks []altToken
	var found bool

	consume := func(regex *regexp.Regexp, typ tokensTy) {
		res := regex.FindAllStringSubmatch(s, 1)
		if res == nil {
			return
		}
		found = true
		r0 := res[0]
		tks = append(tks, altToken{
			typ: typ,
			val: r0[1],
		})
		s = strings.TrimPrefix(s, r0[0])
		s = strings.TrimPrefix(s, " ")
	}

	for i := 0; i < 20; i++ {
		found = false

		consume(empty, tkEmpty)
		consume(literal, tkLiteral)
		consume(singleRune, tkSingleton)
		consume(simpleHex, tkSingleton)
		consume(tenHex, tkSingleton)
		consume(regdot, tkDot)
		consume(regminus, tkMinus)
		consume(regReg, tkRegex)

		if isEmptyOrComment(s) && found {
			return tks, nil
		}

		if !found {
			return nil, fmt.Errorf("couldn't tokenize alternatives")
		}
	}

	return nil, fmt.Errorf("alternative too big")
}

func (tk *altToken) convertRune() {
	if tk.typ != tkSingleton {
		return
	}
	s := tk.val
	tk.val = ""

	rn := []rune(s)
	if len(rn) == 1 {
		tk.r = rn[0]
		return
	}
	num, _ := strconv.ParseInt(s, 16, 0)
	tk.r = rune(num)
}
