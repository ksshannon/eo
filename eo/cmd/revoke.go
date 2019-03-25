// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

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
		// The names are *always* laundered so they are always the
		// same, no need for a lookup
		// "revokee_id",
		"full_revoke_comment",
		"partial_revoke_comment",
		"political",
	})

	eos, err := eo.ParseAllOrders("./data")
	if err != nil {
		log.Fatal(err)
	}

	// Create an index to look up EO by number
	m := map[string]eo.ExecOrder{}
	for _, eo := range eos {
		eon := fmt.Sprintf("%d%s", eo.Number, eo.Suffix)
		// check for duplicates (suffix?)
		if _, ok := m[eon]; ok {
			log.Printf("duplicate eo entry: %+v", eo)
		}
		m[eon] = eo
	}
	for _, eo := range eos {
		revokes := eo.RevokeStrings(true)
		if len(revokes) < 1 {
			continue
		}
		notes := []string{}
		for k, v := range eo.Notes {
			notes = append(notes, fmt.Sprintf("%s: %s", k, v))
		}
		for _, r := range revokes {
			r = r[len("EO "):]
			_, ok := m[r]
			if !ok {
				log.Printf("eo %d can't find revoked %s", eo.Number, r)
			}
			cout.Write([]string{
				fmt.Sprintf("%d%s", eo.Number, eo.Suffix),
				eo.Signed.Format("2006-01-02"),
				eo.Whom(),
				fmt.Sprintf("%s", r),
				m[r].Whom(),
				strings.Join(notes, ";"),
				"skip",
				"-2", // TBD
			})
		}
	}
	cout.Flush()
}
