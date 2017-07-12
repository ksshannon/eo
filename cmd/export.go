// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/ksshannon/eo"
)

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
	})

	eos, err := eo.ParseAllOrders("./data")
	if err != nil {
		panic(err)
	}

	var revokeeEO eo.ExecOrder
	for _, e := range eos {
		w, _ := e.Whom()
		rs := e.RevokeStrings(true)
		for _, revoke := range rs {
			revokeeEO.Number = revoke
			revokee, rid := revokeeEO.Whom()
			cout.Write([]string{
				e.Number,
				e.Notes["Signed"],
				e.Title,
				w,
				revoke,
				revokee,
				fmt.Sprintf("%d", rid),
				e.Notes["Revokes"],
				e.Notes["Revokes in part"],
			})
		}
	}
	cout.Flush()
}
