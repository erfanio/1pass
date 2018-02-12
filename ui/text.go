package ui

import (
	"github.com/therecipe/qt/core"
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
	w.ConnectEnterEvent(func(event *core.QEvent) {
		w.button.DisconnectPaintEvent()
		w.button.Repaint()
	})
	w.ConnectLeaveEvent(func(event *core.QEvent) {
		w.button.ConnectPaintEvent(ignorePaints)
		w.button.Repaint()
	})
}

func (w *TextField) setText(value string) {
	w.text = makeLabel(value, false)
	w.layout.AddWidget(w.text, 0, 0)

	w.button = widgets.NewQPushButton2("Copy", nil)
	w.button.ConnectClicked(func(checked bool) {
		App.Clipboard().SetText(value, gui.QClipboard__Clipboard)
	})
	w.button.ConnectPaintEvent(ignorePaints)

	rect := w.button.FontMetrics().BoundingRect2("Copy")
	w.button.SetFixedSize(core.NewQSize2(rect.Width()+8, rect.Height()+2))

	w.layout.AddWidget(w.button, 0, 0)
}

func ignorePaints(e *gui.QPaintEvent) {}
