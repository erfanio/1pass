package ui

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"os"
)

var App *UI

// UI is a "class" for the app, holds an instance of SearchUI and LoginUI
type UI struct {
	widgets.QApplication

	Search *SearchUI
	Login  *LoginUI
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
}
