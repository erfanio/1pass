package ui

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"os"
)

var App *UI

// UI is a "class" for the app, holds an instance of SearchUI and LoginUI
type UI struct {
	widgets.QApplication

	Search  *SearchUI
	Login   *LoginUI
	Details *DetailsUI
}

func SetupUI() {
	// NewGui will create a new QApplication
	App = NewUI(len(os.Args), os.Args)

	// setup the windows
	App.Search = NewSearchUI(nil, core.Qt__Tool|
		core.Qt__FramelessWindowHint|
		core.Qt__WindowCloseButtonHint|
		core.Qt__WindowStaysOnTopHint)
	App.Login = NewLoginUI(nil, core.Qt__Tool|core.Qt__WindowStaysOnTopHint)
	App.Details = NewDetailsUI(nil, 0)
}

func makeLabel(text string, bold bool) *widgets.QLabel {
	label := widgets.NewQLabel2(text, nil, core.Qt__Widget)
	label.SetTextInteractionFlags(core.Qt__TextSelectableByMouse)
	if bold {
		font := gui.NewQFont()
		font.SetPointSize(14)
		font.SetBold(true)
		label.SetFont(font)
	}
	return label
}
