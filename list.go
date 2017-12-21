package main

import (
	"github.com/therecipe/qt/core"
	"strings"
)

var (
	unfiltered []partialItem
	filtered   []partialItem
	model      *core.QAbstractListModel
)

func populateList(items []partialItem) {
	unfiltered = items
	setupModel()
}

func filter(text string) {
	// data is changing
	model.LayoutAboutToBeChanged(nil, core.QAbstractItemModel__NoLayoutChangeHint)
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

	// data was changed
	model.LayoutChanged(nil, core.QAbstractItemModel__NoLayoutChangeHint)
	ui.Search.UpdateSize()
}

func setupModel() {
	model = core.NewQAbstractListModel(nil)
	model.ConnectFlags(func(index *core.QModelIndex) core.Qt__ItemFlag {
		return core.Qt__ItemIsSelectable | core.Qt__ItemIsEnabled
	})
	model.ConnectData(func(index *core.QModelIndex, role int) *core.QVariant {
		row := index.Row()
		// out of range
		if row >= len(filtered) {
			return core.NewQVariant()
		}

		item := filtered[row]
		if role == int(core.Qt__DisplayRole) {
			return core.NewQVariant17(item.Overview.Title)
		}
		return core.NewQVariant()
	})
	model.ConnectRowCount(func(parent *core.QModelIndex) int {
		return len(filtered)
	})

	ui.Search.List.SetModel(model)
}
