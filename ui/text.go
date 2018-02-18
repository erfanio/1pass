package ui

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type TextField struct {
	widgets.QWidget

	_ func()             `constructor:"init"`
	_ func(string, bool) `slot:"SetText"`

	layout       *widgets.QHBoxLayout
	text         *widgets.QLabel
	copyButton   *widgets.QPushButton
	revealButton *widgets.QPushButton
}

func (w *TextField) init() {
	w.layout = widgets.NewQHBoxLayout()
	w.layout.SetContentsMargins(0, 0, 0, 0)
	w.SetLayout(w.layout)

	w.ConnectSetText(w.setText)
	w.ConnectEnterEvent(func(event *core.QEvent) {
		w.copyButton.DisconnectPaintEvent()
		w.copyButton.Repaint()

		w.revealButton.DisconnectPaintEvent()
		w.revealButton.Repaint()
	})
	w.ConnectLeaveEvent(func(event *core.QEvent) {
		w.copyButton.ConnectPaintEvent(ignorePaints)
		w.copyButton.Repaint()

		w.revealButton.ConnectPaintEvent(ignorePaints)
		w.revealButton.Repaint()
	})
}

func (w *TextField) setText(value string, hidden bool) {
	displayValue := value
	// show if hidden and not empty
	if hidden && len(value) > 0 {
		displayValue = "••••••••••"
	}
	w.text = makeLabel(displayValue, false)
	w.layout.AddWidget(w.text, 0, 0)

	w.copyButton = widgets.NewQPushButton2("Copy", nil)
	w.copyButton.ConnectClicked(func(checked bool) {
		App.Clipboard().SetText(value, gui.QClipboard__Clipboard)
	})
	w.copyButton.ConnectPaintEvent(ignorePaints)

	rect := w.copyButton.FontMetrics().BoundingRect2("Copy")
	w.copyButton.SetFixedSize(core.NewQSize2(rect.Width()+8, rect.Height()+2))

	w.revealButton = widgets.NewQPushButton2("Reveal", nil)
	w.revealButton.ConnectClicked(func(checked bool) {
		w.text.SetText(value)
		w.revealButton.Hide()
	})
	w.revealButton.ConnectPaintEvent(ignorePaints)

	rect = w.revealButton.FontMetrics().BoundingRect2("Reveal")
	w.revealButton.SetFixedSize(core.NewQSize2(rect.Width()+8, rect.Height()+2))

	if !hidden {
		w.revealButton.Hide()
	}
	// don't show the buttons if empty
	if len(value) == 0 {
		w.copyButton.Hide()
		w.revealButton.Hide()
	}

	w.layout.AddWidget(w.copyButton, 0, 0)
	w.layout.AddWidget(w.revealButton, 0, 0)
}

func ignorePaints(e *gui.QPaintEvent) {}
