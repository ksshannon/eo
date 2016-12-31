package eo

import (
	"regexp"
	"strconv"
	"strings"
)

type ExecOrder struct {
	Number string
	Title  string
	Notes  map[string]string
}

var eoMatch = regexp.MustCompile(`EO [0-9]{4,5}`)

func (e *ExecOrder) Revokes() []int {
	var n []int

	s := e.Notes["Revokes"]
	tokens := strings.Split(s, ";")
	for _, t := range tokens {
		if m := eoMatch.FindString(t); m != "" {
			eon, err := strconv.Atoi(m[len("EO "):])
			if err == nil {
				n = append(n, eon)
			}
		}
	}
	return n
}
