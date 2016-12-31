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

func TestParse1937(t *testing.T) {
	fin, _ := os.Open("data/1937.txt")
	defer fin.Close()
	eos := ParseExecOrders(fin)
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
	fin, _ := os.Open("data/1983.txt")
	defer fin.Close()
	eos := ParseExecOrders(fin)
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
	fin, _ := os.Open("data/1979.txt")
	defer fin.Close()
	eos := ParseExecOrders(fin)
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

func TestWhom(t *testing.T) {
	tests := []struct {
		name string
		num  int
	}{
		{
			"Barack Obama",
			13500,
		},
		{
			"Harry S. Truman",
			9540,
		},
		{
			"Ronald Reagan",
			12300,
		},
	}
	var e ExecOrder
	for _, test := range tests {
		e.Number = fmt.Sprintf("%d", test.num)
		if e.Whom() != test.name {
			t.Errorf("invalid whom, expected %s, got %s for %d", e.Whom(), test.name, test.num)
		}
	}
}

// Just attempt to parse all files to weasel out data issues
func TestParseAll(t *testing.T) {

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
	}
}
