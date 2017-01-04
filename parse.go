// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

const delimiter = "Executive Order"

func ParseExecOrders(r io.Reader) []ExecOrder {
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
	if len(eos) < 1 {
		return nil
	}
	return eos[1:]
}
