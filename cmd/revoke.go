package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/ksshannon/eo"
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
		panic(err)
	}

	m := make(map[string]revokeCounts)

	for _, e := range eos {
		w, _ := e.Whom()
		if w == "" {
			fmt.Printf("%+v\n", e)
		}
		who := m[w]
		who.total++
		revoked := e.Revokes()
		who.revoker += len(revoked)
		m[w] = who
		for _, r := range revoked {
			eo := eo.ExecOrder{Number: fmt.Sprintf("%d", r)}
			w, _ := eo.Whom()
			revokee := m[w]
			revokee.revokee++
			m[w] = revokee
		}
	}
	for k, v := range m {
		cout.Write([]string{
			k,
			fmt.Sprintf("%d", v.revoker),
			fmt.Sprintf("%d", v.revokee),
			fmt.Sprintf("%d", v.total),
		})
	}
	cout.Flush()
}
