package main

import (
	"github.com/erfanio/1pass/ui"
	"github.com/therecipe/qt/core"
	"strings"
)

var (
	unfiltered []partialItem
	filtered   []partialItem
)

func setupList() {
	ui.App.Search.SetTextChanged(filter)
	ui.App.Search.SetListDataProviders(listData, listCount, fetchAndCopy, fetchAndOpen)
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
	if row >= len(filtered) || role != int(core.Qt__DisplayRole) {
		return ""
	}
	return filtered[row].Overview.Title
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

// fetchAndOpen
func fetchAndOpen(row int) {
	fetchDetails(filtered[row].UUID, func(item completeItem) {
		values := make(map[string]map[string]string)

		const info = "Info"
		values[info] = make(map[string]string)
		values[info]["Additional Info"] = item.Overview.AdditionalInfo
		values[info]["Title"] = item.Overview.Title
		values[info]["URL"] = item.Overview.URL

		// main fields
		const main = "Main"
		values[main] = make(map[string]string)
		for _, f := range item.Details.Fields {
			values[main][f.Name] = f.Value
		}

		// all fields in all sections
		for _, s := range item.Details.Sections {
			values[s.Name] = make(map[string]string)
			for _, f := range s.Fields {
				values[s.Name][f.Name] = f.Value
			}
		}

		ui.App.Details.Start(item.Overview.Title, values)
	})
}
