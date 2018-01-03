// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"testing"
)

const updateTestData = false

func TestReadFRData(t *testing.T) {
	eos, err := ParseFedRegData(updateTestData)
	if err != nil {
		t.Fatal(err)
	}
	for _, eo := range eos {
		_ = eo
		//fmt.Println(eo.President, eo.Number, eo.RevokeStrings(false))
	}
}
