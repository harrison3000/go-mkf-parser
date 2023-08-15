// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import (
	"fmt"
	"strings"
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

	//TODO what if the end of strings isn't reached?

	return n, nil
}

func (pe *parseEnviroment) matchRule(rule string, input string) (*Node, bool) {
	pe.depth++
	defer func() {
		pe.depth--
	}()

	r := pe.parser.byName[rule]

	var ret *Node

	for _, v := range r.alternatives {
		n, ok := pe.tryAlternative(v, input)
		if !ok {
			continue
		}
		if ret != nil && len(ret.val) > len(n.val) {
			continue
		}
		n.rule = rule
		ret = n
	}

	if ret == nil {
		if r.allowEmpty {
			//TODO improve?
			return &Node{
				rule: rule,
			}, true
		}
		return nil, false
	}

	return ret, true
}

func (pe *parseEnviroment) tryAlternative(alt alternative, input string) (*Node, bool) {
	vLen := 0
	var kids []*Node

	push := func(n *Node) {
		vLen += len(n.val)
		kids = append(kids, n)
	}

	for _, v := range alt.itens {
		s := input[vLen:]

		switch v.kind {
		case itemRune, itemSimpleRange, itemComplexRange:
			c, l := utf8.DecodeRuneInString(s)
			//TODO what about the error?
			var ok bool
			if v.kind == itemComplexRange {
				ok = v.complexRange.inRange(c)
			} else {
				ok = v.runes.inRange(c)
			}
			if !ok {
				return nil, false
			}

			push(&Node{
				val: s[:l],
			})

		case itemRule:
			n, ok := pe.matchRule(v.lit, s)
			if !ok {
				return nil, false
			}
			push(n)

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

			push(&Node{
				val: val,
			})

		case itemLiteral:
			ok := strings.HasPrefix(s, v.lit)
			if !ok {
				return nil, false
			}
			push(&Node{
				val: v.lit,
			})

		default:
			panic("item kind not implemented")
		}
	}

	val := input[:vLen]
	return &Node{
		childs: kids,
		val:    val,
	}, true
}
