// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import "testing"

func TestReadFRData(t *testing.T) {
	_, err := ParseFedRegData(false)
	if err != nil {
		t.Fatal(err)
	}
}
