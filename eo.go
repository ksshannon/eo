// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"regexp"
	"strconv"
	"strings"
)

type ExecOrder struct {
	Number string
	Title  string
	Notes  map[string]string
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

var eoMatch = regexp.MustCompile(`[0-9]+`)

func (e *ExecOrder) Whom() (string, int) {
	m := eoMatch.FindString(e.Number)
	if m == "" {
		return m, -1
	}
	n, err := strconv.Atoi(m)
	if err != nil {
		return "", -1
	}
	return whom(n)
}

var revokeMatch = regexp.MustCompile(`EO [0-9]+`)

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
