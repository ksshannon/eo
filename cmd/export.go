// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-yaml/yaml"
	"github.com/ksshannon/mc/eo"
)

func main() {
	format := flag.String("f", "json", "format(csv,json,yaml)")
	pretty := flag.Bool("p", false, "pretty print")
	flag.Parse()
	fout := os.Stdout
	eos, err := eo.ParseAllOrders("./data")
	if err != nil {
		log.Fatal(err)
	}
	fr, err := eo.FetchFedRegAfterEO(eos[len(eos)-1].Number)
	if err != nil {
		log.Fatal(err)
	}

	eos = append(eos, fr...)

	var b []byte

	switch *format {
	case "csv":
		cout := csv.NewWriter(fout)
		cout.Write([]string{
			"number",
			"suffix",
			"notes",
			"title",
			"president",
			"signed",
		})
		for _, eo := range eos {
			notes := ""
			for k, v := range eo.Notes {
				notes += k + ":" + v + ";"
			}
			cout.Write([]string{
				fmt.Sprintf("%d", eo.Number),
				eo.Suffix,
				notes,
				eo.Title,
				eo.President,
				eo.Signed.Format("01-02-2006"),
			})
		}
	case "json":
		if *pretty {
			b, err = json.MarshalIndent(eos, "", "  ")
		} else {
			b, err = json.Marshal(eos)
		}
	case "yaml":
		b, err = yaml.Marshal(eos)
	}
	if err != nil {
		log.Fatal(err)
	}
	fout.Write(b)
}
