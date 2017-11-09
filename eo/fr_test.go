// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"strings"
	"testing"
)

func TestFetch(t *testing.T) {
	t.Skip()
	eos, err := FetchCurrent()
	if err != nil {
		t.Fatal(err)
	}

	n := eos[0].Number

	for _, eo := range eos[1:] {
		eon := eo.Number
		if eon <= n {
			t.Errorf("%s > previous (%d)", eo.Number, n)
			t.Logf("%+v", eo)
		}
		n = eon
	}
}

func TestFetchCurrent(t *testing.T) {
	eos, err := fetchCurrentFedReg()
	if err != nil {
		t.Fatal(err)
	}
	eo := eos[0]
	// Trump's first was 13765
	if eo.Number != 13765 {
		t.Errorf("fetched wrong eo: %+v", eo)
	}
}

func TestFetchAllOrders(t *testing.T) {
	t.Skip()
	_, err := FetchAllOrders()
	if err != nil {
		t.Error(err)
	}
	oldEO := ParseExecOrdersIn(1998)
	if oldEO == nil {
		t.Error("failed to parse")
	}
	// Computer Software Piracy
	// 13103
	var golden ExecOrder
	for _, eo := range oldEO {
		if eo.Number == 13103 {
			golden = eo
			break
		}
	}
	newEO, err := FetchAllOrders()
	if err != nil {
		t.Error(err)
	}
	for _, eo := range newEO {
		if eo.Number == 13103 && eo.Title != golden.Title {
			t.Errorf("%+v != %+v", golden, eo)
			break
		}
	}
}

func TestOldAgainstNew(t *testing.T) {
	eo := starts[len(starts)-1].start
	old := ParseExecOrdersIn(2017)
	if old == nil {
		t.Fatal("no orders")
	}
	eos, err := fetchFedRegAfterEO(eo - 1)
	if err != nil {
		t.Fatal(err)
	}

	var a, b ExecOrder
	for _, e := range old {
		if e.Number == eo {
			a = e
			break
		}
	}

	for _, e := range eos {
		if e.Number == eo {
			b = e
			break
		}
	}
	// Some titles are title cased in the new data, but not the old
	if a.Number != b.Number || strings.ToLower(a.Title) != strings.ToLower(b.Title) {
		t.Errorf("invalid match: %+v != %+v", a, b)
	}
}
