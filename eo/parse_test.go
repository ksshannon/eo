// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"os"
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
	t.Skip("not needed, looking for conflicts")
	eos, err := ParseAllOrders("./data")
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, e := range eos {
		n := 0
		for k := range e.Notes {
			if strings.Contains(strings.ToLower(k), "revoke") {
				n++
			}
		}
		if n > 1 {
			total++
		}
	}
	if total > 0 {
		t.Logf("%d orders have multiple revoke clauses", total)
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
			if w := e.Whom(); w != "Franklin D. Roosevelt" {
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
		w := e.Whom()
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
		{
			s: `Amends: EO 13043, April 16, 1997; EO 13231, October 16, 2001; EO 13515, October 14, 2009; EO 13538, April 19, 2010; EO 13600, February 9, 2012;
 Supersedes in part: EO 13585, September 30, 2011; EO 13591, November 23, 2011;  EO 13708, September 30, 2015
 Continues: EO 11145, March 7, 1964; EO 11183, October 3, 1964; EO 11287, June 28, 1966; EO 11612, July 26, 1971; EO 12131, May 4, 1979; EO 12216, June 19, 1980; EO 12367, June 15, 1982; EO 12382, September 13, 1982; EO 12829, January 6, 1993; EO 12905, March 25, 1994; EO 12994, March 21, 1996; EO 13231, October 16, 2001; EO 13265, June 6, 2002; EO 13515, October 14, 2009; EO 13521, November 24, 2009; EO 13522, December 9, 2009; EO 13532, February 26, 2010; EO 13538, April 19, 2010; EO 13539, April 21, 2010; EO 13540, April 26, 2010; EO 13549, August 18, 2010; EO 13600, February 9, 2012; EO 13621, July 26, 2012; EO 13631, December 7, 2012; EO 13634, December 21, 2012; EO 13640, January 5, 2013;
 ee: EO 13498, February 5, 2009; EO 13544, June 10, 2010; EO 13555, October 19, 2010`,
			m: map[string]string{
				"Amends":             "EO 13043, April 16, 1997; EO 13231, October 16, 2001; EO 13515, October 14, 2009; EO 13538, April 19, 2010; EO 13600, February 9, 2012;",
				"Supersedes in part": "EO 13585, September 30, 2011; EO 13591, November 23, 2011;  EO 13708, September 30, 2015",
				"Continues":          "EO 11145, March 7, 1964; EO 11183, October 3, 1964; EO 11287, June 28, 1966; EO 11612, July 26, 1971; EO 12131, May 4, 1979; EO 12216, June 19, 1980; EO 12367, June 15, 1982; EO 12382, September 13, 1982; EO 12829, January 6, 1993; EO 12905, March 25, 1994; EO 12994, March 21, 1996; EO 13231, October 16, 2001; EO 13265, June 6, 2002; EO 13515, October 14, 2009; EO 13521, November 24, 2009; EO 13522, December 9, 2009; EO 13532, February 26, 2010; EO 13538, April 19, 2010; EO 13539, April 21, 2010; EO 13540, April 26, 2010; EO 13549, August 18, 2010; EO 13600, February 9, 2012; EO 13621, July 26, 2012; EO 13631, December 7, 2012; EO 13634, December 21, 2012; EO 13640, January 5, 2013;",
				// This is a typo on the data, we are assumming it is See:
				// A fix is in the parseFRNotes as a special case.
				// Original data:
				// "ee": "EO 13498, February 5, 2009; EO 13544, June 10, 2010; EO 13555, October 19, 2010",
				// Fixed:
				"See": "EO 13498, February 5, 2009; EO 13544, June 10, 2010; EO 13555, October 19, 2010",
			},
		},
		{
			s: "Revoked (in part) by: EO 13316, September 17, 2003",
			m: map[string]string{
				"Revoked (in part) by": "EO 13316, September 17, 2003",
			},
		},
	}
	for _, test := range tests {
		m := parseFRNotes(test.s)
		for k, v := range m {
			if test.m[k] != v {
				t.Errorf("got:\n%s:%s, want:\n%s:%s\n", k, v, k, test.m[k])
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
