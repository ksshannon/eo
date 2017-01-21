// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var fedRegURL = url.URL{
	Scheme: "https",
	Host:   "www.federalregister.gov",
	Path:   "/api/v1/documents.json",
}

type fedRegEO struct {
	Number    int    `json:"executive_order_number"`
	Notes     string `json:"executive_order_notes"`
	President struct {
		Name       string `json:"name"`
		Identifier string `json:"identifier"`
	} `json:"president"`
	Significant bool   `json:"significant"`
	SignDate    string `json:"signing_date"`
	Title       string `json:"title"`
}

type fedRegResp struct {
	Count      int        `json:"count"`
	Desc       string     `json:"description"`
	TotalPages int        `json:"total_pages"`
	NextPage   string     `json:"next_page_url"`
	Results    []fedRegEO `json:"results"`
}

// The last EO we have archive data for
const lastEOID = 13738

func FetchCurrent() ([]ExecOrder, error) {
	u := fedRegURL
	q := url.Values{}
	fields := []string{
		"executive_order_notes",
		"executive_order_number",
		"president",
		"significant",
		"signing_date",
		"title",
	}
	for _, f := range fields {
		q.Add("fields[]", f)
	}
	q.Add("per_page", "1000")
	q.Add("conditions[agencies][]", "executive-office-of-the-president")
	q.Add("conditions[type][]", "PRESDOCU")
	q.Add("conditions[presidential_document_type][]", "executive_order")
	//q.Add("order", "relevance")
	q.Add("order", "executive_order_number")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	var results fedRegResp
	if err = json.NewDecoder(resp.Body).Decode(&results); err != nil {
		panic(err)
	}
	var eos []ExecOrder
	var eo ExecOrder
	for _, res := range results.Results {
		if res.Number < 1 {
			continue
		}
		eo.Title = res.Title
		eo.Number = fmt.Sprintf("%d", res.Number)
		eo.Notes = make(map[string]string)
		notes := strings.Split(res.Notes, "\n")
		for _, notes := range notes {
			tokens := strings.Split(strings.TrimSpace(notes), ":")
			if len(tokens) < 1 {
				eo.Notes[tokens[0]] = strings.Join(tokens[1:], ":")
			}
		}
		eos = append(eos, eo)
	}
	return eos, nil
}
