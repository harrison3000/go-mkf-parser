// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

func newComplexRange(base runeRange, excludes []runeRange) *complexRange {
	if !base.valid() {
		return nil
	}

	for _, v := range excludes {
		if !v.valid() {
			return nil
		}
	}

	return &complexRange{
		base:     base,
		excludes: excludes,
	}
}

func (c *complexRange) inRange(char rune) bool {
	if !c.base.inRange(char) {
		return false
	}
	for _, v := range c.excludes {
		if v.inRange(char) {
			return false
		}
	}
	return true
}

func (r runeRange) valid() bool {
	return r[0] <= r[1] && r[0] > 0
}

func (r runeRange) inRange(char rune) bool {
	return r[0] <= char && char <= r[1]
}

// match only exists to implement the matcher interface
func (c *complexRange) match(*parseEnviroment, string) (*Node, bool) {
	panic("shouldn't be here")
}
