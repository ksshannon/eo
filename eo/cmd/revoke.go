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

type revokeCounts struct {
	total   int
	revoker int
	revokee int
}

func main() {
	fout := os.Stdout
	cout := csv.NewWriter(fout)
	cout.Write([]string{
		"president",
		"revoker",
		"revokee",
		"total",
	})

	eos, err := eo.ParseAllOrders("./data")
	if err != nil {
		log.Fatal(err)
	}

	m := make(map[string]revokeCounts)

	for _, e := range eos {
		w := e.Whom()
		if w == "Unknown" {
			fmt.Printf("%+v\n", e)
		}
		who := m[w]
		who.total++
		revoked := e.Revokes()
		who.revoker += len(revoked)
		m[w] = who
		for _, r := range revoked {
			eo := eo.ExecOrder{Number: r}
			w := eo.Whom()
			revokee := m[w]
			revokee.revokee++
			m[w] = revokee
		}
	}
	var ordered = []string{
		eo.Hoover,
		eo.Roosevelt,
		eo.Truman,
		eo.Eisenhower,
		eo.Kennedy,
		eo.Johnson,
		eo.Nixon,
		eo.Ford,
		eo.Carter,
		eo.Reagan,
		eo.BushHW,
		eo.Clinton,
		eo.BushW,
		eo.Obama,
		eo.Trump,
	}

	for _, k := range ordered {
		v := m[k]
		cout.Write([]string{
			k,
			fmt.Sprintf("%d", v.revoker),
			fmt.Sprintf("%d", v.revokee),
			fmt.Sprintf("%d", v.total),
		})
	}
	cout.Flush()
}
