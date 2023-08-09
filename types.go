// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

type Parser struct {
	rules []rule
}

type rule struct {
	name string
	alt  []alternative
}

type alternative []item

type item struct {
	typ matcherType
	lit string
	r   rune
}

type matcherType int8

const (
	mtInvalid matcherType = iota
	mtLiteral
	mtEmpty
	mtRune
	mtSimpleRange
	mtRegex
)
