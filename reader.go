// Copyright 2023 - Harrison Ferreira. All rights reserved.

// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mkf

import (
	"io"
	"strings"
)

type strReader struct {
	original string
	lStack   []int64
	strings.Reader
}

func newStringReader(s string) *strReader {
	ret := &strReader{
		original: s,
	}
	ret.Reset(s)
	return ret
}

func (sr *strReader) PushPos() {
	loc := sr.pos()
	sr.lStack = append(sr.lStack, loc)
}

func (sr *strReader) PopPos() {
	last := len(sr.lStack) - 1
	offset := sr.lStack[last]
	sr.lStack = sr.lStack[:last]
	sr.Seek(offset, io.SeekStart)
}

func (sr *strReader) pos() int64 {
	loc, _ := sr.Seek(0, io.SeekCurrent)
	return loc
}

func (sr *strReader) GetStr() string {
	return sr.original[sr.pos():]
}
