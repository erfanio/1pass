package main

import (
	"github.com/erfanio/1pass/ui"
	"strings"
)

var (
	unfiltered  []partialItem
	filtered    []partialItem
	detailsItem completeItem
)

func setupList() {
	ui.App.Search.SetTextChanged(filter)
	ui.App.Search.SetListDataProviders(listData, listCount, fetchAndCopy, fetchAndOpen)
	ui.App.Details.SetDataProvider(detailsData)
}

func populateList(items []partialItem) {
	unfiltered = items
	ui.App.Search.EnableAndFocus()
}

// filter will search the list of items for matches that show up in the list
func filter(text string) {
	ui.App.Search.ListDataWillChange()
	filtered = make([]partialItem, 0, len(unfiltered))

	// don't bother filtering if text is empty, just hide the list
	if text != "" {
		// filter items
		target := strings.ToLower(text)

		for _, item := range unfiltered {
			if strings.Contains(strings.ToLower(item.Overview.Title), target) ||
				strings.Contains(strings.ToLower(item.Overview.URL), target) ||
				strings.Contains(strings.ToLower(item.Overview.AdditionalInfo), target) {
				filtered = append(filtered, item)
			}
		}
	}

	ui.App.Search.ListDataDidChange()
	ui.App.Search.UpdateSize()
}

// listData will return the data for a row and role
func listData(row, role int) string {
	if row >= len(filtered) {
		return ""
	}

	if role == ui.TitleRole {
		return filtered[row].Overview.Title
	}
	if role == ui.SubtitleRole {
		return filtered[row].Overview.AdditionalInfo
	}
	return ""
}

// listCount will return the number of rows needed to display items
func listCount() int {
	return len(filtered)
}

// fetchAndCopy
func fetchAndCopy(row int) {
	fetchDetails(filtered[row].UUID, func(item completeItem) {
		for _, f := range item.Details.Fields {
			if f.Name == "password" {
				ui.App.Search.Copy(f.Value)
				return
			}
		}
	})
}

// detailsData will supply data to details page
func detailsData() ui.DetailsItem {
	// maybe set later
	username := ""
	password := ""

	sections := make([]ui.DetailsSection, 0)
	detailsFields := make([]ui.DetailsField, 0)

	// details fields
	for _, f := range detailsItem.Details.Fields {
		if f.Name == "username" {
			username = f.Value
		} else if f.Name == "password" {
			password = f.Value
		} else {
			detailsFields = append(detailsFields, ui.DetailsField{
				Title: f.Name,
				Value: f.Value,
			})
		}
	}

	// all fields in all sections
	for _, s := range detailsItem.Details.Sections {
		fields := make([]ui.DetailsField, 0)
		for _, f := range s.Fields {
			fields = append(fields, ui.DetailsField{
				Title: f.Title,
				Value: f.Value,
			})
		}
		sections = append(sections, ui.DetailsSection{
			Title:  s.Title,
			Fields: fields,
		})
	}

	return ui.DetailsItem{
		Title:    detailsItem.Overview.Title,
		URL:      detailsItem.Overview.URL,
		Notes:    detailsItem.Details.Notes,
		Username: username,
		Password: password,
		Fields:   detailsFields,
		Sections: sections,
	}
}

// fetchAndOpen
func fetchAndOpen(row int) {
	fetchDetails(filtered[row].UUID, func(item completeItem) {
		// open and details will fetch data from getData
		detailsItem = item
		ui.App.Details.Start(item.Overview.Title)
	})
}
