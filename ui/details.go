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

type DetailsUI struct {
	widgets.QDialog

	_ func()                                     `constructor:"init"`
	_ func(string, map[string]map[string]string) `slot:"Start"`

	layout *widgets.QGridLayout
}

func (w *DetailsUI) init() {
	w.layout = widgets.NewQGridLayout(nil)
	w.SetLayout(w.layout)

	w.ConnectStart(w.start)
}

func (w *DetailsUI) clear() {
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

func (w *DetailsUI) start(title string, data map[string]map[string]string) {
	// remove all the old widgets
	for w.layout.Count() > 0 {
		item := w.layout.TakeAt(0)
		item.Widget().DestroyQWidget()
		item.DestroyQLayoutItem()
	}

	// set the title
	w.SetWindowTitle(fmt.Sprintf(titleTemplate, title))

	for section, fields := range data {
		w.layout.AddWidget(
			makeLabel(section, true),
			w.layout.RowCount(),
			1,
			core.Qt__AlignLeft,
		)

		for field, value := range fields {
			row := w.layout.RowCount()
			w.layout.AddWidget(
				makeLabel(field, false),
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
	}

	w.Open()
}
