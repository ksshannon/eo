// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/ksshannon/mc/eo"
)

// Desired output:
//
// eo - executive order where the e.o. revocation originates
//
// signed - date issued
//
// title - title of order
//
// president - who revoked
//
// revokes - which executive order the new order revokes
//
// revokee - which presidents' order is being revoked
//
// revokee id - a numeric code for each president
//
// full_revoke_comment - the text from archive.gov
//
// partial_revoke_comment - if the eo partially revokes a past eo (I am only
// looking at full revokes)
//
// political party - indicates which party's order is being revoked. 0 is GOP,
// 1 is Dem. -1 is for orders that were revoked before 1937.

// Numeric ID's for presidents
var pids = map[string]string{
	eo.Unknown:    "-1",
	eo.Hoover:     "31",
	eo.Roosevelt:  "32",
	eo.Truman:     "33",
	eo.Eisenhower: "34",
	eo.Kennedy:    "35",
	eo.Johnson:    "36",
	eo.Nixon:      "37",
	eo.Ford:       "38",
	eo.Carter:     "39",
	eo.Reagan:     "40",
	eo.BushHW:     "41",
	eo.Clinton:    "42",
	eo.BushW:      "43",
	eo.Obama:      "44",
	eo.Trump:      "45",
}

func main() {

	fout := os.Stdout
	cout := csv.NewWriter(fout)

	cout.Write([]string{
		"eo",
		"signed",
		"title",
		"president",
		"revokes",
		"revokee",
		"revokee_id",
		"full_revoke_comment",
		"partial_revoke_comment",
		"political party",
		"all_notes",
	})

	eos, err := eo.ParseAllOrders("./data")
	if err != nil {
		log.Fatal(err)
	}

	// build an index for easier lookup
	m := map[string]eo.ExecOrder{}
	for _, eo := range eos {
		m[eo.Number] = eo
	}

	for _, eo := range eos {
		rvs := eo.Revokes()
		notes := ""
		for k, v := range eo.Notes {
			notes += fmt.Sprintf("%s: %s; ", k, v)
		}
		for _, rv := range rvs {
			r := m[rv]
			cout.Write([]string{
				eo.Number,
				eo.Signed.Format("2006-01-02"),
				eo.Title,
				eo.Whom(),
				r.Number,
				r.Whom(),
				pids[r.Whom()],
				eo.Notes["Revokes"],
				"",
				"",
				notes,
			})
		}
	}
	cout.Flush()
}
