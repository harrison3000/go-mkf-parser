// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import (
	"fmt"
	"unicode/utf8"
)

//TODO save lines to have a kind of "code coverage" for the grammar

func (p *Parser) ParseString(s string) (*Node, error) {
	if len(p.rules) == 0 {
		return nil, fmt.Errorf("empty grammar")
	}

	pe := parseEnviroment{
		parser: p,
	}

	root := p.rules[p.root]

	n, ok := pe.matchRule(root.name, s)
	if !ok {
		return nil, fmt.Errorf("input doesn't match grammar")
	}

	return n, nil
}

func (pe *parseEnviroment) matchRule(rule string, input string) (*Node, bool) {
	pe.depth++
	defer func() {
		pe.depth--
	}()

	ruleidx := pe.parser.mrules[rule] // TODO check ok, just to be sure
	r := &pe.parser.rules[ruleidx]

	matched := make([]*Node, 0, 8)

	for _, v := range r.alt {
		n, ok := pe.tryAlternative(v, input)
		if !ok {
			continue
		}
		n.rule = rule
		matched = append(matched, n)
	}

	if len(matched) == 0 {
		if r.allowEmpty {
			//TODO improve?
			return &Node{}, true
		}
		return nil, false
	}
	if len(matched) > 1 {
		panic("what do we do now?")
	}

	return matched[0], matched != nil
}

func (pe *parseEnviroment) tryAlternative(alt alternative, input string) (*Node, bool) {
	vLen := 0
	var kids []*Node

	for _, v := range alt.itens {
		s := input[vLen:]

		switch v.kind {
		case itemRune, itemSimpleRange:
			c, l := utf8.DecodeRuneInString(s)
			//TODO what about the error?
			if !v.runes.inRange(c) {
				return nil, false
			}
			vLen += l
			kids = append(kids, &Node{
				val: s[:l],
			})

		case itemRule:
			n, ok := pe.matchRule(v.lit, s)
			if !ok {
				return nil, false
			}
			vLen += len(n.val)
			kids = append(kids, n)

		case itemRegex:
			res := v.regex.FindStringIndex(s)
			if res == nil {
				return nil, false
			}
			if res[0] != 0 {
				//TODO explain why
				return nil, false
			}

			val := s[:res[1]]
			vLen += res[1]

			kids = append(kids, &Node{
				val: val,
			})

		default:
			panic("eita deu errado")
		}
	}

	val := input[:vLen]
	return &Node{
		childs: kids,
		val:    val,
	}, true
}
