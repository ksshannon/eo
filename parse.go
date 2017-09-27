// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func ParseExecOrdersIn(year int) []ExecOrder {
	fname := filepath.Join(".", "data", fmt.Sprintf("%d.txt", year))
	st, err := os.Stat(fname)
	if err != nil || st.IsDir() {
		return nil
	}
	fin, err := os.Open(fname)
	if err != nil {
		return nil
	}
	defer fin.Close()
	return ParseExecOrders(fin)
}

func ParseAllOrders(path string) ([]ExecOrder, error) {
	dataFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// TODO(kyle): skip files that don't match a year pattern

	var allOrders []ExecOrder
	for _, fname := range dataFiles {
		fin, err := os.Open(filepath.Join(path, fname.Name()))
		if err != nil {
			return nil, err
		}
		defer fin.Close()
		allOrders = append(allOrders, ParseExecOrders(fin)...)
	}
	return allOrders, nil
}

func parseSigned(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	// Check for the strange signings of FDR such as:
	//
	// EO 8467
	// Signed: 5 FR 2468, July 4, 1940
	if strings.Count(s, ",") > 1 {
		s = strings.TrimSpace(s[strings.Index(s, ",")+1:])
	}
	return time.Parse("January 2, 2006", s)
}

const delimiter = "Executive Order"

var delimitRE = regexp.MustCompile(`^Executive Order [0-9]+(-[A-Z])?$`)

func ParseExecOrders(r io.Reader) []ExecOrder {
	var e ExecOrder
	var err error
	e.Notes = make(map[string]string)
	var eos []ExecOrder
	scn := bufio.NewScanner(r)
	for scn.Scan() {
		text := strings.TrimSpace(scn.Text())
		if text == "" || strings.Index(text, "#") == 1 {
			continue
		}
		if delimitRE.MatchString(text) {
			eos = append(eos, e)
			n := eoMatch.FindString(text)
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
				if tokens[0] == "Signed" {
					e.Signed, err = parseSigned(tokens[1])
					if err != nil {
						log.Print(err)
					}
				}
				e.Notes[tokens[0]] = strings.Join(tokens[1:], ":")
			}
		}
	}
	eos = append(eos, e)
	if len(eos) < 1 {
		return nil
	}
	return eos[1:]
}
