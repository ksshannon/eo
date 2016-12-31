package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ksshannon/eo"
)

func main() {
	// Parse everything in the data file

	dataFiles, err := ioutil.ReadDir("./data")

	if err != nil {
		panic(err)
	}

	fout := os.Stdout
	cout := csv.NewWriter(fout)
	cout.Write([]string{
		"eo",
		"title",
		"president",
		"revokes",
	})

	for _, fname := range dataFiles {
		fin, err := os.Open(filepath.Join("data", fname.Name()))
		if err != nil {
			panic(err)
		}
		defer fin.Close()

		eos := eo.ParseExecOrders(fin)
		if eos == nil {
			panic(fmt.Sprintf("failed to parse %s", fname.Name()))
		}

		for _, e := range eos {
			if e.Notes["Revokes"] != "" {
				cout.Write([]string{
					e.Number,
					e.Title,
					e.Whom(),
					e.Notes["Revokes"],
				})
			}
		}
	}
}
