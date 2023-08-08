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
	comment  = regexp.MustCompile(`^\s*\/\/.*$`)
	ruleName = regexp.MustCompile(`^([a-zA-Z_]+)`)
	ident    = regexp.MustCompile(`^ {4}`)
)

// NewParser returns a new parser... or maybe not
// it accepts a grammar in a modified McKeeman Form
func NewParser(grammar string) (*Parser, error) {
	lines := strings.Split(grammar, "\n")

	exists := map[string]bool{} // TODO store the index and put it in Parser
	var curr rule               //current rule

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
				return nil, fmt.Errorf("too much on line %d", k)
			}

			if exists[n] {
				return nil, fmt.Errorf("rule already defined on line %d", k)
			}
			if curr.name != "" {
				push()
			}
			exists[n] = true
			curr.name = n
			continue
		}

		if ident.MatchString(v) {
			if curr.name == "" {
				return nil, fmt.Errorf("alternative without a name on line %d", k)
			}

			alt, err := str2alt(v[4:], len(curr.alt) == 0)
			if err != nil {
				return nil, fmt.Errorf("error parsing alternative at line %d: %w", k, err)
			}
			curr.alt = append(curr.alt, alt)
			continue
		}
		return nil, fmt.Errorf("didn't understand line %d", k)
	}

	push()

	return &Parser{
		rules: rules,
	}, nil
}

// matcher typers
var (
	literal    = regexp.MustCompile(`^"[^"]*"`)
	singleRune = regexp.MustCompile(`^'.'`)
	simpleHex  = regexp.MustCompile(`^'[0-9A-F]{4,5}'`)
	tenHex     = regexp.MustCompile(`^'10[0-9A-F]{4}'`)
)

func str2alt(s string, allowEmpty bool) (alternative, error) {
	hpf := func(pref string) bool {
		return strings.HasPrefix(s, pref)
	}

	if hpf(`""`) {
		if !allowEmpty {
			return nil, fmt.Errorf("empty not allowed here")
		}
		if !isEmptyOrComment(s[2:]) {
			return nil, fmt.Errorf("empty must be alone")
		}
		return alternative{
			matcher{typ: mtEmpty},
		}, nil
	}

	for {
		if isEmptyOrComment(s) {
			break
		}
		if hpf(`""`) {
			return nil, fmt.Errorf("empty not allowed here")
		}

		return nil, nil
	}

	//TODO disallow empty alternative

	return nil, nil
}

func isEmptyOrComment(s string) bool {
	if strings.TrimSpace(s) == "" {
		return true
	}

	return comment.MatchString(s)
}

func consumeRegex(s string, re *regexp.Regexp) (match, rest string, ok bool) {
	m := re.FindAllStringIndex(s, 1)
	if m == nil {
		return
	}
	m0 := m[0]

	match = s[m0[0]:m0[1]]
	rest = s[m0[1]:]
	ok = true

	return
}
