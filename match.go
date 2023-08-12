// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import (
	"fmt"
)

//TODO save lines to have a kind of "code coverage" for the grammar

func (p *Parser) ParseString(s string) (*Node, error) {
	if len(p.rules) == 0 {
		return nil, fmt.Errorf("empty grammar")
	}

	pe := parseEnviroment{
		parser: p,
		input:  newStringReader(s),
	}

	root := p.rules[p.root]

	n, ok := pe.matchRule(root.name)
	if !ok {
		return nil, fmt.Errorf("input doesn't match grammar")
	}

	return n, nil
}

func (pe *parseEnviroment) matchRule(rule string) (*Node, bool) {
	pe.depth++
	defer func() {
		pe.depth--
	}()

	ruleidx := pe.parser.mrules[rule] // TODO check ok, just to be sure
	r := &pe.parser.rules[ruleidx]

	for _, v := range r.alt {
		n, ok := pe.tryAlternative(v)
		if !ok {
			continue
		}
		n.rule = rule
		return n, true
	}

	if r.allowEmpty {
		//TODO improve
		return &Node{}, true
	}

	return nil, false
}

func (pe *parseEnviroment) tryAlternative(alt alternative) (*Node, bool) {
	pe.input.PushPos()

	fail := pe.input.PopPos

	rest := pe.input.GetStr()
	vLen := 0
	var kids []*Node

	for _, v := range alt.itens {
		switch v.kind {
		case itemRune, itemSimpleRange:
			c, l, _ := pe.input.ReadRune()
			//TODO what about the error?
			if !v.runes.inRange(c) {
				fail()
				return nil, false
			}
			vLen += l
			kids = append(kids, &Node{
				val: string(c), //TODO remove this alloc
			})

		case itemRule:
			n, ok := pe.matchRule(v.lit)
			if !ok {
				fail()
				return nil, false
			}
			vLen += len(n.val)
			kids = append(kids, n)

		case itemRegex:
			str := pe.input.GetStr()
			res := v.regex.FindStringIndex(str)
			if res == nil {
				fail()
				return nil, false
			}
			if res[0] != 0 {
				//TODO explain why
				fail()
				return nil, false
			}

			val := str[:res[1]]
			vLen += res[1]
			pe.input.skip(res[1])

			kids = append(kids, &Node{
				val: val,
			})

		default:
			panic("eita deu errado")
		}
	}

	val := rest[:vLen]
	return &Node{
		childs: kids,
		val:    val,
	}, true
}
