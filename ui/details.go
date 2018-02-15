package ui

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

const DetailsStyles = `
QLabel {
	font-size: 12px;
}

QPushButton {
	font-size: 8px;
	border: 1px solid #cbcccd;
	border-radius: 2px;
	background: #ffffff;
}

QPushButton:pressed {
	border: 1px solid #999999;
	background: #cbcccd;
}
`
const (
	titleTemplate = "Details for %v"
)

type DetailsItem struct {
	Title    string
	URL      string
	Username string
	Password string
	Notes    string
	Sections []DetailsSection
	Fields   []DetailsField
}

type DetailsSection struct {
	Title  string
	Fields []DetailsField
}

type DetailsField struct {
	Title  string
	Value  string
	Hidden bool
}

type DetailsUI struct {
	widgets.QDialog

	_ func()       `constructor:"init"`
	_ func(string) `slot:"Start"`

	layout   *widgets.QGridLayout
	itemData func() DetailsItem
}

func (w *DetailsUI) init() {
	w.layout = widgets.NewQGridLayout(nil)
	w.SetLayout(w.layout)
	w.SetStyleSheet(DetailsStyles)

	w.ConnectStart(w.start)
}

func (w *DetailsUI) SetDataProvider(f func() DetailsItem) {
	w.itemData = f
}

func (w *DetailsUI) makeText(value string, hidden bool) *TextField {
	text := NewTextField(w, 0)
	text.SetText(value, hidden)
	return text
}

func (w *DetailsUI) createRow(key, value string, hidden bool) {
	row := w.layout.RowCount()
	w.layout.AddWidget(
		makeLabel(key, false),
		row,
		0,
		core.Qt__AlignRight|core.Qt__AlignTop,
	)
	w.layout.AddWidget(
		w.makeText(value, hidden),
		row,
		1,
		core.Qt__AlignLeft|core.Qt__AlignTop,
	)
}

func (w *DetailsUI) start(title string) {
	// remove all the old widgets
	for w.layout.Count() > 0 {
		item := w.layout.TakeAt(0)
		item.Widget().DestroyQWidget()
		item.DestroyQLayoutItem()
	}

	// set the title
	w.SetWindowTitle(fmt.Sprintf(titleTemplate, title))

	data := w.itemData()

	w.createRow("Title", data.Title, false)
	w.createRow("Username", data.Username, false)
	w.createRow("Password", data.Password, true)
	w.createRow("URL", data.URL, false)
	w.createRow("Notes", data.Notes, false)
	for _, field := range data.Fields {
		w.createRow(field.Title, field.Value, field.Hidden)
	}

	for _, section := range data.Sections {
		// get rid of related items
		if len(section.Fields) > 0 {
			w.layout.AddWidget(
				makeLabel(section.Title, true),
				w.layout.RowCount(),
				1,
				core.Qt__AlignLeft,
			)

			for _, field := range section.Fields {
				w.createRow(field.Title, field.Value, field.Hidden)
			}
		}
	}

	w.Open()
}
