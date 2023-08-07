// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import "testing"

func TestBasic(t *testing.T) {
	p, e := NewParser(`
//just a test
test
    "wow"

tost
    'a'
`)
	if e != nil {
		t.Fatal("Error:", e)
	}

	_ = p
}
