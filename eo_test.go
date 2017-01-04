// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"fmt"
	"testing"
)

func TestWhom(t *testing.T) {
	tests := []struct {
		name string
		i    int
	}{{"Unknown", 1},
		{"Franklin D. Roosevelt", 6071},
		{"Franklin D. Roosevelt", 9537},
		{"Harry S. Truman", 9540},
		{"Ronald Reagan", 12300},
		{"Barack Obama", 13489}, // FIXME(kyle): failing
		{"Barack Obama", 13490},
		{"Barack Obama", 13500},
		{"Barack Obama", 20000},
	}
	var e ExecOrder
	for _, test := range tests {
		w, _ := whom(test.i)
		if w != test.name {
			t.Errorf("failed whom: %+v (got %s)", test, w)
		}
		e.Number = fmt.Sprintf("%d", test.i)
		w, _ = e.Whom()
		if w != test.name {
			t.Errorf("failed Whom: %+v (got %s)", test, w)
		}
	}
}

func TestWhomAlpha(t *testing.T) {
	e := ExecOrder{Number: "6071-A"}
	w, _ := e.Whom()
	if w != "Franklin D. Roosevelt" {
		t.Errorf("whom failed with alpha")
	}
}
