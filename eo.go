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

var starts = []struct {
	whom  string
	start int
}{
	{"Franklin D. Roosevelt", 6071},
	{"Harry S. Truman", 9538},
	{"Dwight D. Eisenhower", 10432},
	{"John F. Kennedy", 10914},
	{"Lyndon B. Johnson", 11128},
	{"Richard Nixon", 11452},
	{"Gerald R. Ford", 11798},
	{"Jimmy Carter", 11967},
	{"Ronald Reagan", 12287},
	{"George H. W. Bush", 12668},
	{"Bill Clinton", 12834},
	{"George W. Bush", 13198},
	{"Barack Obama", 13489},
}

func whom(order int) (string, int) {
	if order < 6071 {
		return "", -1
	}
	var i int
	for i = 1; i < len(starts); i++ {
		if starts[i].start >= order {
			return starts[i-1].whom, i
		}
	}
	return starts[len(starts)-1].whom, i
}

func (e *ExecOrder) Whom() (string, int) {
	n, err := strconv.Atoi(e.Number)
	if err != nil {
		return "", -1
	}
	return whom(n)
}

var eoMatch = regexp.MustCompile(`EO [0-9]+`)

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
