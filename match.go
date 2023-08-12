// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import (
	"fmt"
)

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

	len := 0
	var kids []*Node

	for _, v := range alt.itens {
		switch v.typ {
		case itemRune:
			c, l, e := pe.input.ReadRune()
			if e != nil || c != v.runes[0] {
				fail()
				return nil, false
			}
			len += l
			kids = append(kids, &Node{
				val: string(c),
			})

		case itemRule:
			n, ok := pe.matchRule(v.lit)
			if !ok {
				fail()
				return nil, false
			}
			//TODO increase len
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

			len += res[1]
			pe.input.skip(res[1])

			kids = append(kids, &Node{
				val: str[:res[1]],
			})

		default:
			panic("eita deu errado")
		}
	}

	return &Node{
		childs: kids,
	}, true
}
