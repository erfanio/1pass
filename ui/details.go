package ui

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

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
	Title string
	Value string
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

	w.ConnectStart(w.start)
}

func (w *DetailsUI) SetDataProvider(f func() DetailsItem) {
	w.itemData = f
}

func makeLabel(text string, bold bool) *widgets.QLabel {
	label := widgets.NewQLabel2(text, nil, core.Qt__Widget)
	label.SetTextInteractionFlags(core.Qt__TextSelectableByMouse)
	if bold {
		font := gui.NewQFont()
		font.SetPixelSize(14)
		font.SetBold(true)
		label.SetFont(font)
	}
	return label
}

func (w *DetailsUI) createRow(key, value string) {
	row := w.layout.RowCount()
	w.layout.AddWidget(
		makeLabel(key, false),
		row,
		0,
		core.Qt__AlignRight|core.Qt__AlignTop,
	)
	w.layout.AddWidget(
		makeLabel(value, false),
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

	w.createRow("Title", data.Title)
	w.createRow("Username", data.Username)
	w.createRow("Password", data.Password)
	w.createRow("URL", data.URL)
	w.createRow("Notes", data.Notes)
	for _, field := range data.Fields {
		w.createRow(field.Title, field.Value)
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
				w.createRow(field.Title, field.Value)
			}
		}
	}

	w.Open()
}
