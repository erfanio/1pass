package main

import (
	"encoding/json"
	"log"
)

type partialItem struct {
	UUID         string   `json:"uuid"`
	TemplateUUID string   `json:"templateUuid"`
	Overview     overview `json:"overview"`
}

type completeItem struct {
	UUID         string   `json:"uuid"`
	TemplateUUID string   `json:"templateUuid"`
	Overview     overview `json:"overview"`
	Details      details  `json:"details"`
}

type overview struct {
	AdditionalInfo string `json:"ainfo"`
	Title          string `json:"title"`
	URL            string `json:"url"`
}

type details struct {
	Fields   []field   `json:"fields"`
	Notes    string    `json:"notesPlain"`
	Sections []section `json:"sections"`
}

type field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type section struct {
	Name   string         `json:"name"`
	Title  string         `json:"title"`
	Fields []sectionField `json:"fields"`
}

type sectionField struct {
	Type  string `json:"k"`
	Title string `json:"t"`
	Value string `json:"v"`
}

func json2list(rawData []byte) (items []partialItem) {
	err := json.Unmarshal(rawData, &items)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func json2item(rawData []byte) (item completeItem) {
	err := json.Unmarshal(rawData, &item)
	if err != nil {
		log.Fatal(err)
	}
	return
}
