// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Source: https://www.archives.gov/federal-register/executive-orders

type ExecOrder2 struct {
	Number    int               `json:"executive_order_number"`
	Suffix    string            `json:"executive_order_number_suffix"`
	Notes     map[string]string `json:"executive_order_notes"`
	Title     string            `json:"title"`
	President string            `json:"president"`
}

func (eo ExecOrder2) Signed() (time.Time, error) {
	s, ok := eo.Notes["Signed"]
	if !ok {
		return time.Time{}, fmt.Errorf("no signed value in notes")
	}
	return time.Parse("January 2, 2006", s)
}

func (eo ExecOrder2) String() string {
	s := fmt.Sprintf("Executive Order %d%s\n", eo.Number, eo.Suffix)
	s += eo.Title + "\n\n"
	//TODO(kyle): order keys
	for k, v := range eo.Notes {
		s += "    " + fmt.Sprintf("%s: %v\n", k, v)
	}
	return s
}

type ExecOrder struct {
	Number int
	Suffix string
	Title  string
	Notes  map[string]string
	Signed time.Time
}

var starts = []struct {
	whom  string
	start int
}{
	{"Franklin D. Roosevelt", 6071},
	{"Harry S. Truman", 9538},
	{"Dwight D. Eisenhower", 10432},
	{"John F. Kennedy", 10914},
	{"Lyndon B. Johnson", 11128},
	{"Richard Nixon", 11452},
	{"Gerald R. Ford", 11798},
	{"Jimmy Carter", 11967},
	{"Ronald Reagan", 12287},
	{"George H. W. Bush", 12668},
	{"Bill Clinton", 12834},
	{"George W. Bush", 13198},
	{"Barack Obama", 13489},
	{"Donald J. Trump", 13758},
}

func whom(order int) (string, int) {
	if order < 6071 {
		return "Unknown", -1
	}
	var i int
	for i = 1; i < len(starts); i++ {
		if starts[i].start > order {
			return starts[i-1].whom, i
		}
	}
	return starts[len(starts)-1].whom, i
}

var eoMatch = regexp.MustCompile(`([0-9]+)(-?[A-Z])?`)
var revokeMatch = regexp.MustCompile(`EO [0-9]+`)
var numMatch = regexp.MustCompile(`[0-9]+`)

func (e *ExecOrder) Whom() (string, int) {
	return whom(e.Number)
}

func (e *ExecOrder) Revokes() []int {
	var n []int
	s := e.Notes["Revokes"]
	tokens := strings.Split(s, ";")
	for _, t := range tokens {
		if m := revokeMatch.FindString(t); m != "" {
			eon, err := strconv.Atoi(m[len("EO "):])
			if err == nil {
				n = append(n, eon)
			}
		}
	}
	return n
}

var revokeStringMatch = regexp.MustCompile(`EO [0-9]+(-[A-Z])?`)

func (e *ExecOrder) RevokeStrings(ignorePartial bool) []string {
	var s []string
	tokens := strings.Split(e.Notes["Revokes"], ";")
	for _, t := range tokens {
		if m := revokeStringMatch.FindString(t); m != "" {
			if ignorePartial && strings.Index(t, "in part") >= 0 {
				continue
			}
			s = append(s, m)
		}
	}
	return s
}
