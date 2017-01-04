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
	if len(allOrders) < 1 {
		t.Fatal("failed to parse")
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
