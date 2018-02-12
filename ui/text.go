package ui

import (
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type TextField struct {
	widgets.QWidget

	_ func()       `constructor:"init"`
	_ func(string) `slot:"SetText"`

	layout *widgets.QHBoxLayout
	text   *widgets.QLabel
	button *widgets.QPushButton
}

func (w *TextField) init() {
	w.layout = widgets.NewQHBoxLayout()
	w.layout.SetContentsMargins(0, 0, 0, 0)
	w.SetLayout(w.layout)

	w.ConnectSetText(w.setText)
}

func (w *TextField) setText(value string) {
	w.text = makeLabel(value, false)
	w.layout.AddWidget(w.text, 0, 0)

	w.button = widgets.NewQPushButton2("Copy", nil)
	w.button.ConnectClicked(func(checked bool) {
		App.Clipboard().SetText(value, gui.QClipboard__Clipboard)
	})
	w.layout.AddWidget(w.button, 0, 0)
}
