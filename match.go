// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import (
	"fmt"
	"regexp"
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
	bn := bunchOfNodes{
		in: input,
	}

	for _, v := range alt.itens {
		s := bn.remaining()

		switch v.kind {
		case itemSimpleRuneRange, itemComplexRange:
			n, ok := tryRune(v, s)
			if !ok {
				return nil, false
			}

			bn.push(n)

		case itemComplex:
			n, ok := v.cplx.match(pe, s)
			if !ok {
				return nil, false
			}

			bn.push(n)

		case itemRule:
			n, ok := pe.matchRule(v.lit, s)
			if !ok {
				return nil, false
			}
			bn.push(n)

		case itemLiteral:
			ok := strings.HasPrefix(s, v.lit)
			if !ok {
				return nil, false
			}
			bn.push(&Node{
				val: v.lit,
			})

		default:
			panic("item kind not implemented")
		}
	}

	return bn.result(), true
}

func tryRune(v item, s string) (*Node, bool) {
	//TODO what about the error rune?
	c, l := utf8.DecodeRuneInString(s)

	var ok bool
	if v.kind == itemComplexRange {
		vc := v.cplx.(*complexRange)
		ok = vc.inRange(c)
	} else {
		ok = v.runes.inRange(c)
	}
	if !ok {
		return nil, false
	}

	val := s[:l]
	return &Node{
		val: val,
	}, true
}

func (cr *cplxRegex) match(_ *parseEnviroment, in string) (*Node, bool) {
	r := (*regexp.Regexp)(cr)
	res := r.FindStringIndex(in)
	if res == nil {
		return nil, false
	}
	if res[0] != 0 {
		//regexes are checked to have a ^ anchor
		panic("impossible")
	}

	val := in[:res[1]]

	return &Node{
		val: val,
	}, true
}

func (bn *bunchOfNodes) push(n *Node) {
	bn.nm += len(n.val)
	bn.ns = append(bn.ns, n)
}

func (bn *bunchOfNodes) remaining() string {
	return bn.in[bn.nm:]
}

func (bn *bunchOfNodes) result() *Node {
	val := bn.in[:bn.nm]
	return &Node{
		childs: bn.ns,
		val:    val,
	}
}
