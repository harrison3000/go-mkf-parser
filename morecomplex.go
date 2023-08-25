// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import (
	"fmt"
	"math"
)

func (k *ruleKnot) match(pe *parseEnviroment, input string) (*Node, bool) {
	bn := bunchOfNodes{
		in: input,
	}

	n, ok := pe.matchRule(k.rule[0], input)
	if !ok {
		return nil, false
	}
	bn.push(n)

	panic("not implemented 0")
}

func (r *ruleRange) match(pe *parseEnviroment, input string) (*Node, bool) {
	bn := bunchOfNodes{
		in: input,
	}

	var matched int32
	for i := 0; i < int(r.ran[1]); i++ {
		rem := bn.remaining()
		n, ok := pe.matchRule(r.rule, rem)
		if !ok {
			break
		}
		matched++
		bn.push(n)
	}

	if matched < r.ran[0] {
		return nil, false
	}

	//TODO rule name?
	return bn.result(), true
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
		//TODO  check if [0] <= [1]
	}

	return ret
}
