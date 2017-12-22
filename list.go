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
	setupModel()
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
				strings.Contains(strings.ToLower(item.Overview.Url), target) ||
				strings.Contains(strings.ToLower(item.Overview.AdditionalInfo), target) {
				filtered = append(filtered, item)
			}
		}
	}

	ui.Search.ListDataDidChange()
	ui.Search.UpdateSize()
}

// connect a list model to data and assign to list
func setupModel() {
	// model needs to be created on gui thread
	ui.Search.SetupListModel(
		func(row, role int) string {
			// out of range
			if row >= len(filtered) || role != int(core.Qt__DisplayRole) {
				return ""
			}
			return filtered[row].Overview.Title
		},
		func() int {
			return len(filtered)
		},
	)
	ui.Search.Enable()
}
