// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var revokeTags = []string{
	"Revoked (in part) by",
	"Revoked by",
	"Revoked by in part and supplemented by",
	"Revoked by in part by",
	"Revoked in part and supplemented by",
	"Revoked in part by",
	"Revoked in pary by",
	"Revoked part by",
	"Revokes",
	"Revokes advisory committees established by",
	"Revokes in part",
	"Revokes in part and supplements",
}

var noteTypos = map[string]struct {
	count   int
	correct string
}{
	"Abolished by": {1, ""},
	"Abolishes":    {1, ""},
	"Advisory Council for Minority Enterprise Continued by": {1, ""},
	"Amend by":               {1, "Amended by"},
	"Amended by":             {1126, ""},
	"Amended by (continued)": {4, ""},
	"Amends":                 {1640, ""},
	"Authority repealed by":  {3, ""},
	"Board continued by":     {2, ""},
	"Citizens' Advisory Council on the Status of Women continued by": {1, ""},
	"Commission continued by":                      {4, ""},
	"Commission extended by":                       {1, ""},
	"Committee continued by":                       {8, ""},
	"Committee terminated by":                      {1, ""},
	"Consumer Advisory Council continued by":       {1, ""},
	"Continued by":                                 {43, ""},
	"Continues":                                    {10, ""},
	"Continues advisory committees established by": {2, ""},
	"Continues certain committee established by":   {2, ""},
	"Continues certain committees established by":  {3, ""},
	"Continues committees established by":          {3, ""},
	"Continues in effect":                          {1, ""},
	"Continuing various committees established by": {1, ""},
	"Council abolished by":                         {1, ""},
	"Council continued by":                         {1, ""},
	"Effective":                                    {1, ""},
	"Effective date":                               {1, ""},
	"Expired":                                      {1, ""},
	"Extended by":                                  {1, ""},
	"Federal Advisory Council on Occupational Safety and Health continued by": {1, ""},
	"Federal Register Page and Date":                                          {103, ""},
	"Federal Register correction page and date":                               {3, ""},
	"Federal Register page and date":                                          {4860, ""},
	"Modified by":                                                             {1, ""},
	"Modifies":                                                                {2, ""},
	"NOTE":                                                                    {1, ""},
	"Note":                                                                    {520, ""},
	"Nullified by":                                                            {2, ""},
	"President's Committee on the National Medal of Science continued by": {2, ""},
	"Provisionally supersedes":                   {1, ""},
	"Provisions extended by":                     {1, ""},
	"Ratified by":                                {10, ""},
	"Repealed by":                                {1, ""},
	"Rescinded by":                               {2, ""},
	"Review Board continued by":                  {1, ""},
	"Revoked (in part) by":                       {2, ""},
	"Revoked by":                                 {1572, ""},
	"Revoked by in part and supplemented by":     {1, ""},
	"Revoked by in part by":                      {1, "Revoked in part by"},
	"Revoked in part and supplemented by":        {1, ""},
	"Revoked in part by":                         {31, ""},
	"Revoked in pary by":                         {1, "Revoked in part by"},
	"Revoked part by":                            {1, ""},
	"Revokes":                                    {662, ""},
	"Revokes advisory committees established by": {1, ""},
	"Revokes in part":                            {38, ""},
	"Revokes in part and supplements":            {1, ""},
	"See":                                        {1635, ""},
	"Signed":                                     {5393, ""},
	"Superseded by":                              {439, ""},
	"Superseded in part by":                      {27, ""},
	"Superseded or revoked by":                   {4, ""},
	"Superseded or revoked in part by":           {2, ""},
	"Superseded to extent inconsistent by":       {1, ""},
	"Supersedes":                                 {327, ""},
	"Supersedes (export control provisions)":     {1, ""},
	"Supersedes in part":                         {33, ""},
	"Supersedes or revokes":                      {1, ""},
	"Supplemented by":                            {5, ""},
	"Supplements":                                {4, ""},
	"Suspended by":                               {3, ""},
	"Suspends":                                   {2, ""},
	"Suspersedes":                                {1, ""},
	"Task Force continued by":                    {1, ""},
	"Terminated by":                              {1, ""},
	"Terminates":                                 {1, ""},
	"Terminates committees in":                   {1, ""},
	"Voided by":                                  {1, ""},
	"ee":                                         {1, "See"},
}

// Source: https://www.archives.gov/federal-register/executive-orders

// ExecOrder represents a single order issued by a president
type ExecOrder struct {
	// Executive order number, some have a trailing alpha character such as 'A',
	// or '-A'
	Number    string            `json:"number",yaml:"number"`
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
	s := fmt.Sprintf("Executive Order %s\n", eo.Number)
	s += eo.Title + "\n\n"
	for k, v := range eo.Notes {
		s += "    " + fmt.Sprintf("%s: %v\n", k, v)
	}
	return s
}

var tenures = []struct {
	whom  string
	start int
	end   int
}{
	{Hoover, 5075, 6070}, // No actual data for HH, just the EO #
	{Roosevelt, 6071, 9537},
	{Truman, 9538, 10431},
	{Eisenhower, 10432, 10913},
	{Kennedy, 10914, 11127},
	{Johnson, 11128, 11451},
	{Nixon, 11452, 11797},
	{Ford, 11798, 11966},
	{Carter, 11967, 12286},
	{Reagan, 12287, 12667},
	{BushHW, 12668, 12833},
	{Clinton, 12834, 13197},
	{BushW, 13198, 13488},
	{Obama, 13489, 13764},
	{Trump, 13765, math.MaxInt32},
}

func whom(order int) string {
	for _, s := range tenures {
		if order >= s.start && order <= s.end {
			return s.whom
		}
	}
	return "Unknown"
}

var eoMatch = regexp.MustCompile(`([0-9]+)(-?[A-Z])?`)
var revokeMatch = regexp.MustCompile(`EO [0-9]+`)
var numMatch = regexp.MustCompile(`[0-9]+`)

func (e *ExecOrder) Whom() string {
	if e.Number == "" {
		return Unknown
	}
	return whom(e.AsInt())
}

func (e ExecOrder) AsInt() int {
	s := ""
	for _, c := range e.Number {
		if c < '0' || c > '9' {
			break
		}
		s += string(c)
	}
	x, _ := strconv.Atoi(s)
	return x
}

// Return the order numbers of the orders that an order revokes
//
// TODO(kyle): return a full EO, so we can have the suffix.
func (e *ExecOrder) Revokes() []string {
	var n []string
	s := e.Notes["Revokes"]
	tokens := strings.Split(s, ";")
	for _, t := range tokens {
		if m := revokeMatch.FindString(t); m != "" {
			n = append(n, m[len("EO "):])
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
