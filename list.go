package main

import (
	"encoding/json"
	"github.com/therecipe/qt/core"
	"log"
	"strings"
)

type partial struct {
	Uuid     string   `json:"uuid"`
	Overview overview `json:"overview"`
}

type overview struct {
	AdditionalInfo string `json:"ainfo"`
	Title          string `json:"title"`
	Url            string `json:"url"`
}

var (
	items         []partial
	filteredItems []partial
	model         *core.QAbstractListModel
)

func initList() {
	// setup model
	initModel()
	search.SetListModel(model)

	search.OnTextChanged(func(text string) {
		// don't bother filtering if text is empty, just hide the list
		if text == "" {
			search.HideList()
			return
		}

		// search list for lower case input text
		target := strings.ToLower(text)

		// data is changing
		model.LayoutAboutToBeChanged(nil, core.QAbstractItemModel__NoLayoutChangeHint)

		// two lists, one for full object and one of strings that goes to gui
		filteredItems = make([]partial, 0, len(items))

		for _, item := range items {
			// search in title, url and additional info
			title := strings.ToLower(item.Overview.Title)
			url := strings.ToLower(item.Overview.Url)
			ainfo := strings.ToLower(item.Overview.AdditionalInfo)

			if strings.Contains(title, target) || strings.Contains(url, target) || strings.Contains(ainfo, target) {
				filteredItems = append(filteredItems, item)
			}
		}

		// data was changed
		model.LayoutChanged(nil, core.QAbstractItemModel__NoLayoutChangeHint)
		search.UpdateSize()
	})
}

func initModel() {
	// create a new model and connect necessary functions
	// model uses filteredItems as data source
	model = core.NewQAbstractListModel(nil)
	// flags
	model.ConnectFlags(func(index *core.QModelIndex) core.Qt__ItemFlag {
		return core.Qt__ItemIsSelectable | core.Qt__ItemIsEnabled
	})
	// data
	model.ConnectData(func(index *core.QModelIndex, role int) *core.QVariant {
		row := index.Row()
		// out of range
		if row >= len(filteredItems) {
			return core.NewQVariant()
		}

		item := filteredItems[row]
		if role == int(core.Qt__DisplayRole) {
			return core.NewQVariant17(item.Overview.Title)
		}
		return core.NewQVariant()
	})
	// count
	model.ConnectRowCount(func(parent *core.QModelIndex) int {
		return len(filteredItems)
	})
}

func setupListData(data []byte) {
	// update items from json data
	items = nil
	err := json.Unmarshal(data, &items)
	if err != nil {
		// crash
		log.Fatal(err)
	}
}
