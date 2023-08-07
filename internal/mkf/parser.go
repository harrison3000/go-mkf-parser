// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import (
	"fmt"
	"regexp"
	"strings"
)

type Parser struct {
	rules []rule
}

type rule struct {
	name string
	alt  []matcher
}

type matcher string

var (
	comment  = regexp.MustCompile(`^\s*\/\/.*$`)
	ruleName = regexp.MustCompile(`^([a-zA-Z_]+)\s*(\/\/.*)?$`)
	ident    = regexp.MustCompile(`^ {4}`)
)

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
		if strings.TrimSpace(v) == "" {
			continue
		}

		if comment.MatchString(v) {
			//we just ignore comments for now
			continue
		}

		if n := ruleName.FindStringSubmatch(v); n != nil {
			if exists[n[0]] {
				return nil, fmt.Errorf("rule already defined on line %d", k)
			}
			if curr.name != "" {
				push()
			}
			exists[n[0]] = true
			curr.name = n[0]
			continue
		}

		if ident.MatchString(v) {
			if curr.name == "" {
				return nil, fmt.Errorf("alternative without a name on line %d", k)
			}

			curr.alt = append(curr.alt, str2matcher(v[4:]))
			continue
		}
		return nil, fmt.Errorf("didn't understand line %d", k)
	}

	push()

	return &Parser{
		rules: rules,
	}, nil
}

func str2matcher(s string) matcher {
	//TODO implement
	var m matcher
	return m
}
