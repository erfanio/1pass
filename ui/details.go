package ui

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"strings"
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
	RawJSON  string
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
	if len(strings.TrimSpace(value)) == 0 {
		return
	}
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

func (w *DetailsUI) createSectionRow(section DetailsSection) {
	if len(strings.TrimSpace(section.Title)) > 0 {
		w.layout.AddWidget(
			makeLabel(section.Title, true),
			w.layout.RowCount(),
			1,
			core.Qt__AlignLeft,
		)
	}
}

func (w *DetailsUI) createJSONRow(json string) {
	textEdit := widgets.NewQPlainTextEdit2(json, nil)
	textEdit.Hide()

	button := widgets.NewQPushButton2("Show JSON", nil)
	button.ConnectClicked(func(checked bool) {
		button.Hide()
		textEdit.Show()
	})

	row := w.layout.RowCount()
	w.layout.AddWidget(
		makeLabel("JSON", false),
		row,
		0,
		core.Qt__AlignRight|core.Qt__AlignTop,
	)
	w.layout.AddWidget(
		textEdit,
		row,
		1,
		core.Qt__AlignLeft|core.Qt__AlignTop,
	)
	w.layout.AddWidget(
		button,
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
		if !isSectionEmpty(section) {
			w.createSectionRow(section)

			for _, field := range section.Fields {
				w.createRow(field.Title, field.Value, field.Hidden)
			}
		}
	}

	w.createJSONRow(data.RawJSON)

	w.Open()
}

func isSectionEmpty(section DetailsSection) bool {
	for _, field := range section.Fields {
		if len(strings.TrimSpace(field.Value)) > 0 {
			return false
		}
	}
	return true
}
