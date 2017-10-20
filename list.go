package main

import (
	"encoding/json"
	"github.com/therecipe/qt/core"
	"log"
	"os"
	"os/exec"
	"strings"
)

type partialItem struct {
	Uuid         string   `json:"uuid"`
	TemplateUuid string   `json:"templateUuid"`
	Overview     overview `json:"overview"`
}

type completeItem struct {
	Uuid         string   `json:"uuid"`
	TemplateUuid string   `json:"templateUuid"`
	Overview     overview `json:"overview"`
	Details      details  `json:"details"`
}

type overview struct {
	AdditionalInfo string `json:"ainfo"`
	Title          string `json:"title"`
	Url            string `json:"url"`
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
	Name  string `json:"n"`
	Title string `json:"t"`
	Value string `json:"v"`
}

var (
	items         []partialItem
	filteredItems []partialItem
	model         *core.QAbstractListModel
)

func initList(sessionEnv string) {
	// get items from op
	cmd := exec.Command("/bin/op", "list", "items")
	cmd.Env = append(os.Environ(), sessionEnv)
	out, err := cmd.Output()

	// login if list couldn't be fetched
	if err != nil {
		logExit(err)
		initLogin()
		return
	}

	// update items from json data
	items = nil
	err = json.Unmarshal(out, &items)
	if err != nil {
		// crash
		log.Fatal(err)
	}

	initListGui(sessionEnv)
}

func initListGui(sessionEnv string) {
	// setup model
	initListModel()
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
		filteredItems = make([]partialItem, 0, len(items))

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

	search.ContextMenuData(func(row int) map[string]string {
		// get items from op
		cmd := exec.Command("/bin/op", "get", "item", filteredItems[row].Uuid)
		cmd.Env = append(os.Environ(), sessionEnv)
		out, err := cmd.Output()

		// login if list couldn't be fetched
		if err != nil {
			logExit(err)
			initLogin()
			return nil
		}

		// get item
		var item completeItem
		err = json.Unmarshal(out, &item)
		if err != nil {
			// crash
			log.Fatal(err)
		}

		// basic info +
		results := map[string]string{
			"notes": item.Details.Notes,
			"url":   item.Overview.Url,
			"JSON":  string(out),
		}
		// all the field +
		for _, f := range item.Details.Fields {
			if f.Value != "" {
				results[f.Name] = f.Value
			}
		}
		// all fields in all sections
		for _, s := range item.Details.Sections {
			for _, f := range s.Fields {
				if f.Value != "" {
					results[f.Title] = f.Value
				}
			}
		}

		return results
	})
}

func initListModel() {
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

func logExit(err error) {
	exitErr, ok := err.(*exec.ExitError)
	if ok {
		log.Print(err, string(exitErr.Stderr))
	} else {
		log.Print(err)
	}
}
