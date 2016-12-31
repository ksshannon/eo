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
		"revokee",
		"full_revoke_comment",
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

		var tmpOrder eo.ExecOrder
		for _, e := range eos {
			rev := e.Revokes()
			if rev != nil {
				for _, r := range rev {
					tmpOrder.Number = fmt.Sprintf("%d", r)
					w, _ := e.Whom()
					_, tw := tmpOrder.Whom()
					cout.Write([]string{
						e.Number,
						e.Title,
						w,
						fmt.Sprintf("%d", r),
						fmt.Sprintf("%d", tw),
						e.Notes["Revokes"],
					})
				}
			}
		}
	}
}
