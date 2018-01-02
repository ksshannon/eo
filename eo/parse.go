// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// ParseExecOrdersIn reads the orders from the data folder for the specified
// year.  If the year isn't in the data folder, or any other error is
// encountered, a nil slice is returned.
//
// TODO(kyle): return a valid error on error
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

const (
	delimiter = "Executive Order"

	// The federalregister.gov data starts at this eo number, don't parse past
	// this.
	firstFR = 12893
)

var delimitRE = regexp.MustCompile(`^Executive Order [0-9]+(-?[A-Z])?$`)

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
			e.President, _ = e.Whom()
			eos = append(eos, e)
			// Add support for the suffix extraction
			matches := eoMatch.FindStringSubmatch(text)
			e.Number, err = strconv.Atoi(matches[1])
			if err != nil {
				log.Print(err)
			}
			if e.Number >= firstFR {
				break
			}
			e.Suffix = matches[2]
			if e.Suffix != "" && e.Suffix[0] == '-' {
				e.Suffix = e.Suffix[1:]
			}
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

// ParseAllExecOrders parses all of the internal text orders, then parses the
// internal JSON files.  Finally, it makes a web request for any new orders
func ParseAllExecOrders() []ExecOrder {
	return nil
}

var archiveURL = url.URL{
	Scheme: "https",
	Host:   "www.archives.gov",
	Path:   "/federal-register/executive-orders/{{YEAR}}-{{WHOM}}.html",
}

var archiveResolv = []struct {
	year int
	name string
}{
	{1933, "roosevelt"},
	//{2009, "obama"},
}

// The Executive Orders in the web pages appear like:
//<p><a name="13490"></a> <strong><a class="pdfImage" href="http://www.gpo.gov/fdsys/pkg/FR-2009-01-26/pdf/E9-1719.pdf">Executive Order 13490</a></strong><br />
//  Ethics Commitments by Executive Branch Personnel</p>
//
//<ul>
//  <li>Signed: January 21, 2009</li>
//  <li>Federal Register page and date: 74 FR 4673, January 26, 2009</li>
//  <li>Superseded by: <a href="/federal-register/executive-orders/2017-trump#13770">EO 13770</a>, January 28, 2017</li>
//</ul>
//
//<hr />
//
// The anchor holds the EO number It appears we can grab it from the name
// attribute.
//
func parseWeb() ([]ExecOrder, error) {
	var eos []ExecOrder
	for _, ar := range archiveResolv {
		u := archiveURL
		u.Path = strings.Replace(u.Path, "{{YEAR}}", strconv.Itoa(ar.year), 1)
		u.Path = strings.Replace(u.Path, "{{WHOM}}", ar.name, 1)
		resp, err := http.Get(u.String())
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		t := html.NewTokenizer(resp.Body)
		if t == nil {
			return nil, io.EOF
		}
		var eo ExecOrder
		for tkn := t.Next(); tkn != html.ErrorToken; tkn = t.Next() {
			if tkn != html.TextToken {
				continue
			}
			text := string(bytes.TrimSpace(t.Text()))
			// This is the title of the EO
			if delimitRE.MatchString(text) {
				eo.Notes = map[string]string{}
				eo.Number, err = strconv.Atoi(strings.Fields(text)[2])
				if err != nil {
					return nil, err
				}
				t.Next()
				eo.Title = string(t.Text())
				tkn = t.Next()
				text = string(t.Text())
				eo.Notes = map[string]string{}
				for tkn != html.ErrorToken && !delimitRE.MatchString(text) {
					fmt.Println("XXX", text)
					k := strings.Split(text, ":")
					if len(k) == 2 {
						eo.Notes[k[0]] = k[1]
					}
					tkn = t.Next()
					text = string(t.Text())
				}
				eos = append(eos, eo)
			}
		}
	}
	return eos, nil
}
