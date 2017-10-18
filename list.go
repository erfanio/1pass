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
	model         *core.QStringListModel
)

func initList() {
	// setup model
	model = core.NewQStringListModel2([]string{}, nil)
	search.SetListModel(model)

	search.OnTextChanged(func(text string) {
		// search list for lower case input text
		target := strings.ToLower(text)

		// two lists, one for full object and one of strings that goes to gui
		stringItems := make([]string, 0, len(items))
		filteredItems = make([]partial, 0, len(items))

		for _, item := range items {
			// search in title, url and additional info
			title := strings.ToLower(item.Overview.Title)
			url := strings.ToLower(item.Overview.Url)
			ainfo := strings.ToLower(item.Overview.AdditionalInfo)

			if strings.Contains(title, target) || strings.Contains(url, target) || strings.Contains(ainfo, target) {
				stringItems = append(stringItems, title)
				filteredItems = append(filteredItems, item)
			}
		}

		// put found items in the gui list
		model.SetStringList(stringItems)
		search.UpdateSize()
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
