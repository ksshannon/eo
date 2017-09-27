// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
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

type fedRegResp2 struct {
	Count       int64  `json:"count"`
	Description string `json:"description"`
	Results     []struct {
		BodyHTMLURL          string `json:"body_html_url"`
		Citation             string `json:"citation"`
		DocumentNumber       string `json:"document_number"`
		EndPage              int64  `json:"end_page"`
		ExecutiveOrderNotes  string `json:"executive_order_notes"`
		ExecutiveOrderNumber int64  `json:"executive_order_number"`
		FullTextXMLURL       string `json:"full_text_xml_url"`
		HTMLURL              string `json:"html_url"`
		JSONURL              string `json:"json_url"`
		PdfURL               string `json:"pdf_url"`
		PublicationDate      string `json:"publication_date"`
		SigningDate          string `json:"signing_date"`
		StartPage            int64  `json:"start_page"`
		Title                string `json:"title"`
	} `json:"results"`
	TotalPages int64 `json:"total_pages"`
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

func FetchAllOrders() ([]ExecOrder, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "federalregister.gov",
		Path:   "/api/v1/documents.json",
	}
	q := url.Values{}
	q.Add("conditions[correction]", "0")
	q.Add("conditions[presidential_document_type_id]", "2")
	q.Add("conditions[type]", "PRESDOCU")
	q.Add("fields[]", "citation")
	q.Add("fields[]", "document_number")
	q.Add("fields[]", "end_page")
	q.Add("fields[]", "executive_order_notes")
	q.Add("fields[]", "executive_order_number")
	q.Add("fields[]", "html_url")
	q.Add("fields[]", "pdf_url")
	q.Add("fields[]", "publication_date")
	q.Add("fields[]", "signing_date")
	q.Add("fields[]", "start_page")
	q.Add("fields[]", "title")
	q.Add("fields[]", "full_text_xml_url")
	q.Add("fields[]", "body_html_url")
	q.Add("fields[]", "json_url")
	q.Add("order", "executive_order_number")
	q.Add("per_page", "1000")
	u.RawQuery = q.Encode()
	_ = u
	/*
		resp, err := http.Get(u.String())
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	*/
	pres := []string{"clinton", "bushw", "obama", "trump"}
	var buf []byte
	for _, p := range pres {
		path := filepath.Join("data", "fr", p) + ".json"
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		buf = append(buf, b...)
	}
	var result fedRegResp2
	//err = json.NewDecoder(resp.Body).Decode(&result)
	err := json.NewDecoder(bytes.NewReader(buf)).Decode(&result)
	if err != nil {
		return nil, err
	}
	var eos []ExecOrder
	for _, r := range result.Results {
		eo := ExecOrder{
			Number: fmt.Sprintf("%d", r.ExecutiveOrderNumber),
			Title:  r.Title,
			Notes:  map[string]string{},
		}
		tokens := strings.Split(r.ExecutiveOrderNotes, ":")
		if len(tokens) > 1 {
			eo.Notes[tokens[0]] = strings.Join(tokens[1:], ":")
		}
		eos = append(eos, eo)
	}
	return eos, nil
}
