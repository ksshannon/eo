// Copyright 2016 Kyle Shannon.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eo

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	Clinton = "william-j-clinton"
	Bush    = "george-w-bush"
	Obama   = "barack-obama"
	Trump   = "donald-trump"
	Current = Trump
)

var fedRegURL = url.URL{
	Scheme: "https",
	Host:   "www.federalregister.gov",
	Path:   "/api/v1/documents.json",
}

type fedRegResp struct {
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

// updateFRData downloads the presidental EOs and saves them in
// data/fr/{{president}}.json.  An error is returned if encountered.  All
// fields are downloaded.
func updateFRData(whom string) error {
	u := fedRegURL
	q := url.Values{}
	q.Add("conditions[correction]", "0")
	q.Add("conditions[presidential_document_type_id]", "2")
	q.Add("conditions[type]", "PRESDOCU")
	q.Add("order", "executive_order_number")
	q.Add("per_page", "1000")
	q.Add("conditions[president]", whom)
	for _, field := range allFedRegFields {
		q.Add("fields[]", field)
	}
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fout, err := os.Create(filepath.Join(".", "data", "fr", whom+".json"))
	if err != nil {
		return err
	}
	defer fout.Close()
	_, err = io.Copy(fout, resp.Body)
	return err
}

func ParseFedRegData(update bool) ([]ExecOrder, error) {
	if update {
		if err := updateFRData(Current); err != nil {
			return nil, err
		}
	}
	return readLocalFedReg()
}

func readLocalFedReg() ([]ExecOrder, error) {
	// We'd like them in EO order:
	pres := []string{
		Clinton,
		Bush,
		Obama,
		Trump,
	}
	var buf []byte
	for _, p := range pres {
		path := filepath.Join(".", "data", "fr", p+".json")
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		buf = append(buf, b...)
	}
	return parseFedRegJSON(bytes.NewReader(buf))
}

func parseFedRegJSON(r io.Reader) ([]ExecOrder, error) {
	var result fedRegResp
	err := json.NewDecoder(r).Decode(&result)
	if err != nil {
		return nil, err
	}
	var eos []ExecOrder
	for _, r := range result.Results {
		w, _ := whom(int(r.ExecutiveOrderNumber))
		eo := ExecOrder{
			Number:    int(r.ExecutiveOrderNumber),
			Title:     r.Title,
			President: w,
			Notes:     map[string]string{},
		}
		tokens := strings.Split(r.ExecutiveOrderNotes, ":")
		if len(tokens) > 1 {
			eo.Notes[tokens[0]] = strings.Join(tokens[1:], ":")
		}
		eo.Signed, err = time.Parse("2006-01-02", r.SigningDate)
		eos = append(eos, eo)
	}
	return eos, nil
}

var allFedRegFields = []string{
	"abstract",
	"action",
	"agencies",
	"agency_names",
	"body_html_url",
	"cfr_references",
	"citation",
	"comment_url",
	"comments_close_on",
	"correction_of",
	"corrections",
	"dates",
	"docket_id",
	"docket_ids",
	"document_number",
	"effective_on",
	"end_page",
	"excerpts",
	"executive_order_notes",
	"executive_order_number",
	"full_text_xml_url",
	"html_url",
	"images",
	"json_url",
	"mods_url",
	"page_length",
	"pdf_url",
	"president",
	"public_inspection_pdf_url",
	"publication_date",
	"raw_text_url",
	"regulation_id_number_info",
	"regulation_id_numbers",
	"regulations_dot_gov_info",
	"regulations_dot_gov_url",
	"significant",
	"signing_date",
	"start_page",
	"subtype",
	"title",
	"toc_doc",
	"toc_subject",
	"topics",
	"type",
	"volume",
}
