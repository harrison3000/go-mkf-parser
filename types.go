// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

type Parser struct {
	rules  []rule
	mrules map[string]int
}

type rule struct {
	name string
	alt  []alternative
}

type alternative []item

type item struct {
	typ itemType
	lit string
	r   rune
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
)

type grammarParseError struct {
	err  string
	line int
}
