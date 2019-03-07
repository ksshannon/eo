// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type frEO struct {
	Count       int64  `json:"count"`
	Description string `json:"description"`
	Results     []struct {
		BodyHTMLURL          string `json:"body_html_url"`
		Citation             string `json:"citation"`
		DispositionNotes     string `json:"disposition_notes"`
		DocumentNumber       string `json:"document_number"`
		EndPage              int64  `json:"end_page"`
		ExecutiveOrderNumber int64  `json:"executive_order_number"`
		FullTextXMLURL       string `json:"full_text_xml_url"`
		HTMLURL              string `json:"html_url"`
		JSONURL              string `json:"json_url"`
		PdfURL               string `json:"pdf_url"`
		PublicationDate      string `json:"publication_date"`
		SigningDate          string `json:"signing_date"`
		StartPage            int64  `json:"start_page"`
		Subtype              string `json:"subtype"`
		Title                string `json:"title"`
		Type                 string `json:"type"`
	} `json:"results"`
	TotalPages int64 `json:"total_pages"`
}

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

var noteMatch = regexp.MustCompile(`[A-Za-z\(\) ]+:`)

func parseFRNotes(s string) map[string]string {
	m := map[string]string{}
	match := noteMatch.FindAllStringIndex(s, -1)
	for i := 0; i < len(match); i++ {
		var a, b int
		a = match[i][1]
		if i < len(match)-1 {
			b = match[i+1][0]
		} else {
			b = len(s)
		}
		key := strings.TrimSpace(s[match[i][0]:match[i][1]])
		key = key[:len(key)-1]
		typo := noteTypos[key]
		if typo.count < 2 && typo.correct != "" {
			key = typo.correct
		}
		m[key] = strings.TrimSpace(s[a:b])
	}
	return m
}

func ParseAllOrders(path string) ([]ExecOrder, error) {
	var files []string
	for i := 1937; i < 1994; i++ {
		files = append(files, fmt.Sprintf("%d.txt", i))
	}
	var allOrders []ExecOrder
	for _, f := range files {
		fin, err := os.Open(filepath.Join(path, f))
		if err != nil {
			return nil, err
		}
		allOrders = append(allOrders, ParseExecOrders(fin)...)
		fin.Close()
	}

	fin, err := os.Open(filepath.Join(path, "fr.json"))
	if err != nil {
		return nil, err
	}
	defer fin.Close()
	var fr frEO
	err = json.NewDecoder(fin).Decode(&fr)
	if err != nil {
		return nil, err
	}
	for _, feo := range fr.Results {
		w, n := whom(int(feo.ExecutiveOrderNumber))
		if n < 0 {
			return nil, fmt.Errorf("invalid eo: %d", feo.ExecutiveOrderNumber)
		}
		when, _ := time.Parse("2006-01-02", feo.SigningDate)
		allOrders = append(allOrders, ExecOrder{
			Notes:     parseFRNotes(feo.DispositionNotes),
			Number:    int(feo.ExecutiveOrderNumber),
			President: w,
			Signed:    when,
			Suffix:    "",
			Title:     feo.Title,
		})
	}
	// The FR orders are stored most current first, lets sort them all.
	sort.Slice(allOrders, func(i, j int) bool {
		if allOrders[i].Number == allOrders[j].Number {
			return allOrders[i].Suffix < allOrders[j].Suffix
		}
		return allOrders[i].Number < allOrders[j].Number
	})
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

var delimitRE = regexp.MustCompile(`^Executive Order [0-9]+(-?[A-Z])?$`)

func ParseExecOrders(r io.Reader) []ExecOrder {
	var e ExecOrder
	var err error
	e.Notes = make(map[string]string)
	var eos []ExecOrder
	scn := bufio.NewScanner(r)
	for scn.Scan() {
		text := strings.TrimSpace(scn.Text())
		if text == "" || text[0] == '#' {
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
			for i, t := range tokens {
				tokens[i] = strings.TrimSpace(t)
			}
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
