// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import "regexp"

type Parser struct {
	byName map[string]*rule
	rules  []rule
	root   int
}

type rule struct {
	name         string
	alternatives []alternative
	allowEmpty   bool
}

type alternative struct {
	itens []item
}

type cMatcher interface {
	match(string) (string, bool)
}

type item struct {
	cplx cMatcher

	lit string

	runes runeRange
	kind  itemKind
}

type itemKind int8

const (
	itemInvalid itemKind = iota
	itemLiteral
	itemEmpty
	itemSimpleRuneRange
	itemComplexRange
	itemComplex
	itemRule
)

type grammarParseError struct {
	err  string
	line int
}

type complexRange struct {
	excludes []runeRange
	base     runeRange
}

type runeRange [2]rune

type Node struct {
	rule   string
	val    string
	childs []*Node
}

type parseEnviroment struct {
	parser *Parser
	depth  int
}

type cplxRegex regexp.Regexp

type ruleKnot struct {
	rule [2]string
	char rune
}

type ruleRange struct {
	rule string
	ran  [2]int32
}

type bunchOfNodes struct {
	ns []*Node
	in string
	nm int
}
