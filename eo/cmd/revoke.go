// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
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
		for _, rv := range rvs {
			r := m[rv]
			cout.Write([]string{
				eo.Number,
				eo.Signed.Format("2006-01-02"),
				eo.Title,
				eo.Whom(),
				r.Number,
				r.Whom(),
				"",
				eo.Notes["Revokes"],
				"",
				"",
			})
		}
	}
	cout.Flush()
}
