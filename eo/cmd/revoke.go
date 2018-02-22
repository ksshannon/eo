// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/ksshannon/mc/eo"
)

type revokeCounts struct {
	total   int
	revoker int
	revokee int
}

func main() {
	update := flag.Bool("u", false, "update the local json data before running")
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
		panic(err)
	}

	_ = *update
	/*
		freos, err := eo.ParseFedRegData(*update)
		if err != nil {
			panic(err)
		}
		eos = append(eos, freos...)
	*/

	m := make(map[string]revokeCounts)

	for _, e := range eos {
		w, _ := e.Whom()
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
			w, _ := eo.Whom()
			revokee := m[w]
			revokee.revokee++
			m[w] = revokee
		}
	}
	var ordered = []string{
		"Herbert Hoover",
		"Franklin D. Roosevelt",
		"Harry S. Truman",
		"Dwight D. Eisenhower",
		"John F. Kennedy",
		"Lyndon B. Johnson",
		"Richard Nixon",
		"Gerald R. Ford",
		"Jimmy Carter",
		"Ronald Reagan",
		"George H. W. Bush",
		"Bill Clinton",
		"George W. Bush",
		"Barack Obama",
		"Donald J. Trump",
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
