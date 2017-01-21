// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"strconv"
	"testing"
)

func TestFetch(t *testing.T) {
	t.Skip()
	eos, err := FetchCurrent()
	if err != nil {
		t.Fatal(err)
	}

	n, _ := strconv.Atoi(eos[0].Number)

	for _, eo := range eos[1:] {
		eon, _ := strconv.Atoi(eo.Number)
		if eon <= n {
			t.Errorf("%s > previous (%d)", eo.Number, n)
			t.Logf("%+v", eo)
		}
		n = eon
	}
}
