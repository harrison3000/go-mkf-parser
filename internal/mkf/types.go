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

type alternative []matcher

type matcher struct {
	typ matcherType
}

type matcherType int8

const (
	mtLiteral matcherType = iota
	mtRune
	mtSimpleRange
	mtRegex
)
