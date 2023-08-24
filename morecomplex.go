// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import (
	"fmt"
	"math"
)

func (k *ruleKnot) match(string) (string, bool) {
	panic("not implemented 0")
}

func (r *ruleRange) match(string) (string, bool) {
	panic("not implemented 1")
}

func mkRuleRange(rule string, rg string) *ruleRange {
	ret := &ruleRange{
		rule: rule,
	}
	switch rg {
	case "?":
		ret.ran = [2]int32{0, 1}

	case "+":
		ret.ran = [2]int32{1, math.MaxInt32}

	case "*":
		ret.ran = [2]int32{0, math.MaxInt32}

	default:
		var a, b int32
		n, _ := fmt.Sscanf(rg, "{%d,%d}", &a, &b)
		if n == 1 {
			ret.ran = [2]int32{a, a}
		} else if n == 2 {
			ret.ran = [2]int32{a, b}
		} else {
			panic("???")
		}
	}

	return ret
}
