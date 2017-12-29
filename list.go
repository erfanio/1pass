package main

import (
	"github.com/therecipe/qt/core"
	"strings"
)

var (
	unfiltered []partialItem
	filtered   []partialItem
)

func populateList(items []partialItem) {
	unfiltered = items
	ui.Search.EnableAndFocus()
}

func filter(text string) {
	ui.Search.ListDataWillChange()
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

	ui.Search.ListDataDidChange()
	ui.Search.UpdateSize()
}

// ListData will return the data for a row and role
func ListData(row, role int) string {
	if row >= len(filtered) || role != int(core.Qt__DisplayRole) {
		return ""
	}
	return filtered[row].Overview.Title
}

// ListCount will return the number of rows needed to display items
func ListCount() int {
	return len(filtered)
}
