package eo

import (
	"bufio"
	"io"
	"strings"
)

const delimiter = "Executive Order"

func parseExecOrders(r io.Reader) []ExecOrder {
	var e ExecOrder
	var eos []ExecOrder
	scn := bufio.NewScanner(r)
	for scn.Scan() {
		text := strings.TrimSpace(scn.Text())
		if text == "" {
			continue
		}
		if strings.HasPrefix(text, delimiter) {
			eos = append(eos, e)
			n := strings.TrimSpace(text[len(delimiter):])
			e.Number = n
			e.Title = ""
			e.Notes = make(map[string]string)
			continue
		}
		if e.Title == "" {
			e.Title = text
		} else {
			tokens := strings.Split(text, ":")
			if len(tokens) > 1 {
				e.Notes[tokens[0]] = strings.Join(tokens[1:], ":")
			}
		}
	}
	return eos[1:]
}
