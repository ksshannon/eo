// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestParseInvalidYear(t *testing.T) {
	eos := ParseExecOrdersIn(1900)
	if eos != nil {
		t.Fatalf("opened non-existent year")
	}
}

func TestParse1937(t *testing.T) {
	eos := ParseExecOrdersIn(1937)
	if eos == nil {
		t.Fatal("parsing failed")
	}
	// Check the data in the first order
	e := eos[0]
	if e.Number != "7532" {
		t.Errorf("incorrect number: %s", e.Number)
	}
	if strings.Index(e.Title, "Shinnecock") < 0 {
		t.Errorf("incorrect title: %s", e.Title)
	}
	if len(e.Notes) < 1 {
		t.Fatal("invalid notes")
	}
	if n, ok := e.Notes["Revoked by"]; !ok {
		t.Errorf("invalid notes: %+v", e.Notes)
	} else if strings.Index(n, "Public") < 0 {
		t.Errorf("invalid notes: %+v", e.Notes)
	}
}

func TestParse1983(t *testing.T) {
	eos := ParseExecOrdersIn(1983)
	if eos == nil {
		t.Fatal("parsing failed")
	}

	// Find 12407, it should be revoke 12314
	found := false
	for _, e := range eos {
		if e.Number == "12407" {
			found = true
			if strings.Index(e.Notes["Revokes"], "12314") < 0 {
				t.Errorf("invalid revokes note: %s", e.Notes["Revokes"])
			}
		}
	}
	if !found {
		t.Error("couldn't find proper order (12407)")
	}
}

func TestMultiRevoke(t *testing.T) {
	eos := ParseExecOrdersIn(1979)
	if eos == nil {
		t.Fatal("parsing failed")
	}

	found := false
	// Find 12148, revokes many orders, including 10242
	for _, e := range eos {
		if e.Number == "12148" {
			revokes := e.Revokes()
			for _, n := range revokes {
				if n == 10242 {
					found = true
					break
				}
			}
		}
	}
	if !found {
		t.Error("didn't find 10242 in the revoke notes")
	}
}

// Just attempt to parse all files to weasel out data issues
func TestParseAll(t *testing.T) {
	allOrders, err := ParseAllOrders("./data")
	if err != nil {
		t.Fatal(err)
	}
	var invalid int
	for i, e := range allOrders {
		this, err := strconv.Atoi(e.Number)
		if err == nil && i > 0 {
			last, err := strconv.Atoi(allOrders[i-1].Number)
			if err == nil {
				if this-last != 1 {
					t.Log(last, this)
					invalid++
				}
			}
		}
	}
	// grep -E '^Executive Order [0-9]+(-[A-Z])?$' data/*.txt | wc -l
	// reports 6275.
	const orderCount = 6279
	if len(allOrders) != orderCount {
		t.Errorf("parsed %d orders, expected %d", len(allOrders), orderCount)
	}
	if invalid > 0 {
		t.Logf("possible invalid: %d", invalid)
	}
}

func TestShortEONumber(t *testing.T) {
	fin, _ := os.Open("data/1986.txt")
	defer fin.Close()
	eos := ParseExecOrders(fin)
	if eos == nil {
		t.Fatal("parsing failed")
	}

	found := false
	// Find 12553, revokes many orders, including a short EO number, 723
	for _, e := range eos {
		if e.Number == "12553" {
			revokes := e.Revokes()
			for _, n := range revokes {
				if n == 723 {
					found = true
					break
				}
			}
		}
	}
	if !found {
		t.Error("didn't find 723 in the revoke notes")
	}
}

func TestRevokesInPart(t *testing.T) {
	t.Skip()
	dataFiles, err := ioutil.ReadDir("./data")

	if err != nil {
		t.Fatal(err)
	}

	fout := os.Stdout
	cout := csv.NewWriter(fout)
	cout.Write([]string{
		"eo",
		"title",
		"president",
		"revokes",
	})

	conflict := 0

	for _, fname := range dataFiles {
		fin, err := os.Open(filepath.Join("data", fname.Name()))
		if err != nil {
			panic(err)
		}
		defer fin.Close()

		eos := ParseExecOrders(fin)
		if eos == nil {
			t.Fatal(fmt.Sprintf("failed to parse %s", fname.Name()))
		}
		for _, e := range eos {
			if _, hasInPart := e.Notes["Revokes in part"]; hasInPart == true {
				if strings.Index("in part", strings.ToLower(e.Notes["Revokes"])) >= 0 {
					t.Logf("file: %s eo: %s has revokes(in part): %s and revokes: %s",
						fname.Name(), e.Number, e.Notes["Revokes"], e.Notes["Revokes in part"])
					conflict++
				}
			}
		}
	}
	if conflict > 0 {
		t.Logf("%d conflicts in revoke/revoke in part", conflict)
	}
}

func TestAlphaEO(t *testing.T) {
	// 1037 has 7677-A
	eos := ParseExecOrdersIn(1937)
	if eos == nil {
		t.Fatalf("failed to parse 1937")
	}
	found := false
	for _, e := range eos {
		if e.Number == "7677-A" {
			found = true
			if w, _ := e.Whom(); w != "Franklin D. Roosevelt" {
				t.Errorf("Whom failed on AlphaEO: %s", w)
			}
		}
	}
	if !found {
		t.Error("failed to parse alpha eo 7677-A")
	}
}

// Issues with this specific order
func Test9379(t *testing.T) {
	eos := ParseExecOrdersIn(1943)
	if eos == nil {
		t.Fatal("failed to parse")
	}
	const n = "9379"
	var found bool
	for _, e := range eos {
		if e.Number == n {
			found = true
			if e.Title == "" {
				t.Error("failed to extract title")
			}
			break
		}
	}
	if !found {
		t.Error("failed to find 9379")
	}
}

func TestMisses(t *testing.T) {
	tests := []int{
		8029 + 1,  // 8031
		8315 + 1,  // 8317
		8623 + 1,  // 8625
		8889 + 1,  // 8889
		9004 + 1,  // 9006
		9291 + 1,  // 9293
		9411 + 1,  // 9413
		9507 + 1,  // 9509
		9671 + 1,  // 9673
		9816 + 1,  // 9818
		9917 + 1,  // 9919
		10024 + 1, // 10026
		10093 + 1, // 10095
		10198 + 1, // 10200
		10316 + 1, // 10318
		10510 + 1, // 10512
		10571 + 1, // 10571
		10583 + 1, // 10585
		10648 + 1, // 10650
		10692 + 1, // 10694
		10746 + 1, // 10748
		10796 + 1, // 10798
		10856 + 1, // 10858
		10898 + 1, // 10900
		10982 + 1, // 10984
		11071 + 1, // 11073
		11133 + 1, // 11135
		11189 + 1, // 11191
		11263 + 1, // 11265
		11320 + 1, // 11322
		11385 + 1, // 11387
		11441 + 1, // 11443
		11502 + 1, // 11504
		11574 + 1, // 11576
		11637 + 1, // 11639
		11692 + 1, // 11694
		11756 + 1, // 11758
		11825 + 1, // 11827
		11892 + 1, // 11894
		11948 + 1, // 11950
		12031 + 1, // 12033
		12109 + 1, // 12111
		12186 + 1, // 12188
		12259 + 1, // 12261
		12335 + 1, // 12337
		12398 + 1, // 12400
		12455 + 1, // 12457
		12496 + 1, // 12498
		12541 + 1, // 12543
		12578 + 1, // 12580
		12621 + 1, // 12623
		12661 + 1, // 12663
		12697 + 1, // 12699
		12740 + 1, // 12742
		12786 + 1, // 12788
		12826 + 1, // 12828
		12889 + 1, // 12891
		12943 + 1, // 12945
		12983 + 1, // 12985
		13032 + 1, // 13034
		13070 + 1, // 13072
		13108 + 1, // 13110
		13143 + 1, // 13145
		13184 + 1, // 13186
		13250 + 1, // 13252
		13281 + 1, // 13283
		13322 + 1, // 13324
		13367 + 1, // 13369
		13393 + 1, // 13395
		13420 + 1, // 13422
		13452 + 1, // 13454
		13482 + 1, // 13484
		13526 + 1, // 13528
		13561 + 1, // 13563
		13595 + 1, // 13597
		13634 + 1, // 13636
		13654 + 1, // 13656
		13685 + 1, // 13687
		13714 + 1, // 13716
	}
	eos, err := ParseAllOrders("./data")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range tests {
		found := false
		for _, e := range eos {
			if e.Number == fmt.Sprintf("%d", test) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("failed to find %d", test)
		}
	}
}

func TestCount(t *testing.T) {
	m := make(map[string]int)
	eos, err := ParseAllOrders("./data")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range eos {
		w, _ := e.Whom()
		m[w]++
	}
}

func TestRevokeString(t *testing.T) {
	eos := ParseExecOrdersIn(1940)
	var found bool
	var eo ExecOrder
	for _, e := range eos {
		if e.Number == "8346" {
			eo = e
			found = true
			break
		}
	}
	if !found {
		t.Fatal("failed to find eo")
	}
	s := eo.RevokeStrings(false)
	if s == nil {
		t.Error("failed to parse strings")
	}

	found = false
	for _, r := range s {
		if r == "EO 3653-A" {
			found = true
			break
		}
	}
	if !found {
		t.Error("failed to find EO 3653-A")
	}
	// Make sure we extract the right amount of EOs in int and string versions
	if len(eo.Revokes()) != len(eo.RevokeStrings(false)) {
		t.Error("length of Revokes() and RevokeStrings() mis-match")
	}
}
