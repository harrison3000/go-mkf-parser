// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	ruleName = regexp.MustCompile(`^([a-zA-Z_]+)`)
	ident    = regexp.MustCompile(`^( {4}|\t)`)
)

// NewParser returns a new parser... or maybe not
// it accepts a grammar in a modified McKeeman Form
func NewParser(grammar string) (*Parser, error) {
	lines := strings.Split(grammar, "\n")

	used := map[string]bool{}

	mrules := map[string]int{} // TODO store the index and put it in Parser
	var curr rule              //current rule

	var rules []rule

	push := func() {
		rules = append(rules, curr)
		curr = rule{}
	}

	for k, v := range lines {
		if isEmptyOrComment(v) {
			//we just ignore empty lines and comments for now
			continue
		}

		if n, rest, ok := consumeRegex(v, ruleName); ok {
			if !isEmptyOrComment(rest) {
				return nil, newParseError("unexpected content after rule name", k)
			}

			if _, ok := mrules[n]; ok {
				return nil, newParseError("duplicate rule", k)
			}
			if curr.name != "" {
				push()
			}
			mrules[n] = len(rules)
			curr.name = n
			continue
		}

		if _, _, ok := consumeRegex(v, ident); ok {
			if curr.name == "" {
				return nil, newParseError("orphaned alternative", k)
			}

			alt, err := str2alt(v, len(curr.alt) == 0)
			if err != nil {
				//TODO improve this
				return nil, fmt.Errorf("error parsing alternative at line %d: %w", k, err)
			}

			for _, item := range alt.itens {
				if item.typ == itemRule {
					used[item.lit] = true
				}
			}

			curr.alt = append(curr.alt, alt)
			continue
		}
		return nil, newParseError("unable to parse grammar", k)
	}

	push()

	for k := range used {
		if _, ok := mrules[k]; !ok {
			return nil, fmt.Errorf("rule not found: %s", k)
		}
	}

	//TODO warn unused?

	return &Parser{
		rules:  rules,
		mrules: mrules,
	}, nil
}

func isEmptyOrComment(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}

	return strings.HasPrefix(s, "//")
}

func consumeRegex(s string, re *regexp.Regexp) (groupMatch, rest string, ok bool) {
	m := re.FindAllStringSubmatch(s, 1)
	if m == nil {
		return
	}
	m0 := m[0]

	groupMatch = m0[1]
	rest, ok = strings.CutPrefix(s, m0[0])
	return
}

func newParseError(err string, line int) *grammarParseError {
	return &grammarParseError{
		line: line,
		err:  err,
	}
}

func (g *grammarParseError) Error() string {
	return fmt.Sprintf("%s, on line: %d", g.err, g.line)
}
