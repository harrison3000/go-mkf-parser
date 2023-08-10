// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import "regexp"

type Parser struct {
	mrules map[string]int
	rules  []rule
}

type rule struct {
	name string
	alt  []alternative
}

type alternative struct {
	itens []item
	empty bool
}

type item struct {
	lit string

	complexRange *complexRange
	regex        *regexp.Regexp

	runes runeRange
	typ   itemType
}

type itemType int8

const (
	itemInvalid itemType = iota
	itemLiteral
	itemEmpty
	itemRune
	itemSimpleRange
	itemComplexRange
	itemRegex
	itemRule
)

type grammarParseError struct {
	err  string
	line int
}

type complexRange struct {
	base     runeRange
	excludes []runeRange
}

type runeRange [2]rune

type Result struct {
}
