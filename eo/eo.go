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

// ExecOrder represents a single order issued by a president
type ExecOrder struct {
	Number    int               `json:"number",yaml:"number"`
	Suffix    string            `json:"suffix",yaml:"suffix"`
	Notes     map[string]string `json:"notes",yaml:"notes"`
	Title     string            `json:"title",yaml:"title"`
	President string            `json:"president",yaml:"president"`
	Signed    time.Time         `json:"signed",yaml:"signed"`
}

const (
	Unknown    = "Unknown"
	Hoover     = "Herbert Hoover"
	Roosevelt  = "Franklin D. Roosevelt"
	Truman     = "Harry S. Truman"
	Eisenhower = "Dwight D. Eisenhower"
	Kennedy    = "John F. Kennedy"
	Johnson    = "Lyndon B. Johnson"
	Nixon      = "Richard Nixon"
	Ford       = "Gerald R. Ford"
	Carter     = "Jimmy Carter"
	Reagan     = "Ronald Reagan"
	BushHW     = "George H. W. Bush"
	Clinton    = "Bill Clinton"
	BushW      = "George W. Bush"
	Obama      = "Barack Obama"
	Trump      = "Donald J. Trump"
)

// String returns a formated order that closely matches the format from
// Roosevelt to 1994, when the federal register takes over.
func (eo ExecOrder) String() string {
	s := fmt.Sprintf("Executive Order %d%s\n", eo.Number, eo.Suffix)
	s += eo.Title + "\n\n"
	for k, v := range eo.Notes {
		s += "    " + fmt.Sprintf("%s: %v\n", k, v)
	}
	return s
}

var starts = []struct {
	whom  string
	start int
}{
	{Hoover, 5075}, // No actual data for HH, just the EO #
	{Roosevelt, 6071},
	{Truman, 9538},
	{Eisenhower, 10432},
	{Kennedy, 10914},
	{Johnson, 11128},
	{Nixon, 11452},
	{Ford, 11798},
	{Carter, 11967},
	{Reagan, 12287},
	{BushHW, 12668},
	{Clinton, 12834},
	{BushW, 13198},
	{Obama, 13489},
	{Trump, 13765},
}

func whom(order int) (string, int) {
	if order < starts[0].start {
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

// Return the order numbers of the orders that an order revokes
//
// TODO(kyle): return a full EO, so we can have the suffix.
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
