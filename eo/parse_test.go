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
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseInvalidYear(t *testing.T) {
	eos := ParseExecOrdersIn(1900)
	if eos != nil {
		t.Fatalf("opened non-existent year")
	}
}

func TestEODelimitMatch(t *testing.T) {
	tests := []struct {
		s string
		b bool
	}{
		{"Executive Order 1234", true},
		{"Executive Order 1234A", true},
		{"Executive Order 1234B", true},
		{"Executive Order 1234-A", true},
		{"Executive Order 1234AA", false},
		{"Executive Order 1234-AA", false},
	}
	for _, test := range tests {
		if delimitRE.MatchString(test.s) != test.b {
			t.Errorf("invalid match: %s returned %t", test.s, test.b)
		}
	}
}

func TestEONumberSubmatch(t *testing.T) {
	t.Skip("write me")
	/*
		tests := []struct {
			s string
			m []string
		}{
			{"1234", []string{"1234", "", ""}},
			{"1234A", []string{"1234A", "1234", "A"}},
			{"1234B", []string{"1234B", "1234", "B"}},
			{"1234-A", []string{"1234", "1234", "-A"}},
			{"1234AA", []string{"", "", ""}},
			{"1234-AA", []string{"", "", ""}},
		}
		for _, test := range tests {
			m := delimitRE.SubMatchString(test.s)
			for i, mm := range m {
				if mm != test.m[i] {
					t.Errorf("invalid match: %s returned %t", mm, test.m[i])
				}
			}
		}
	*/
}

func TestParse1937(t *testing.T) {
	eos := ParseExecOrdersIn(1937)
	if eos == nil {
		t.Fatal("parsing failed")
	}
	// Check the data in the first order
	e := eos[0]
	if e.Number != 7532 {
		t.Errorf("incorrect number: %d", e.Number)
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
	if e.President != "Franklin D. Roosevelt" {
		t.Errorf("invalid president: %s", e.President)
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
		if e.Number == 12407 {
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
		if e.Number == 12148 {
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
		if e.Number == 12553 {
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
					t.Logf("file: %s eo: %d has revokes(in part): %s and revokes: %s",
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
		if e.Number == 7677 && e.Suffix == "A" {
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
	const n = 9379
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
		if e.Number == 8346 {
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

func TestSigned(t *testing.T) {
	tests := []struct {
		n int
		s time.Time
	}{
		{7726, time.Date(1937, 10, 12, 0, 0, 0, 0, time.UTC)},
		// Two 'Signed' keys
		{7729, time.Date(1937, 10, 16, 0, 0, 0, 0, time.UTC)},
	}
	eos := ParseExecOrdersIn(1937)
	for _, test := range tests {
		for _, eo := range eos {
			if eo.Number == test.n && eo.Signed != test.s {
				t.Errorf("signed mismatch, exp: %s, got %s", test.s, eo.Signed)
			}
		}
	}
}

func TestFRNoteMatch(t *testing.T) {
	tests := []struct {
		s string
		m map[string]string
	}{
		{
			s: `Supersedes: EO 12869, September 30, 1993; Revokes: EO 12878, November 5, 1993; Superseded by: EO 13062, September 29, 1997; See: EO 12887, December 23, 1993; EO 12912, April 29, 1994`,
			m: map[string]string{
				"Supersedes":    "EO 12869, September 30, 1993;",
				"Revokes":       "EO 12878, November 5, 1993;",
				"Superseded by": "EO 13062, September 29, 1997;",
				"See":           "EO 12887, December 23, 1993; EO 12912, April 29, 1994",
			},
		},
	}
	for _, test := range tests {
		m := parseFRNotes(test.s)
		for k, v := range m {
			if test.m[k] != v {
				t.Errorf("got: %s:%s, want: %s:%s\n", k, v, k, test.m[k])
			}
		}
	}
}

func TestFRData(t *testing.T) {
	tests := []ExecOrder{
		{
			Notes: map[string]string{
				"See":                   "EO 13303, May 22, 2003; EO 13364, November 29, 2004; EO 13438, July 17, 2007; Notice of May 20, 2004; Notice of May 19, 2005; Notice of May 18, 2006; Notice of May 18, 2007; Notice of May 20, 2008; Notice of May 19, 2009; Notice of May 12, 2010; EO 13668, May 27, 2014;",
				"Superseded in part by": "EO 13350, July 29, 2004",
			},
			Number:    13315,
			President: "George W. Bush",
			Signed:    time.Date(2003, 8, 28, 0, 0, 0, 0, time.UTC),
			Suffix:    "",
			Title:     "Blocking Property of the Former Iraqi Regime, Its Senior Officials and Their Family Members, and Taking Certain Other Actions",
		},
		{
			Notes:     map[string]string{},
			Number:    13735,
			President: "Barack Obama",
			Signed:    time.Date(2016, 8, 12, 0, 0, 0, 0, time.UTC),
			Suffix:    "",
			Title:     "Providing an Order of Succession Within the Department of the Treasury",
		},
	}
	eos, err := ParseAllOrders("./data")
	if err != nil {
		t.Fatal(err)
	}
	m := map[int]int{}
	for i, eo := range eos {
		m[eo.Number] = i
	}
	for _, test := range tests {
		i, ok := m[test.Number]
		if !ok {
			t.Errorf("failed to parse eo: %d", test.Number)
		}
		got := eos[i]
		if !reflect.DeepEqual(test, got) {
			t.Errorf("failed to parse fr eo data, got: %+v, want: %+v", got, test)
		}
	}
}
