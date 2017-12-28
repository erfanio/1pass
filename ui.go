package main

import (
	"github.com/therecipe/qt/widgets"
	"os"
)

var ui *UI

// UI is a "class" for the app, holds an instance of SearchUI and LoginUI
type UI struct {
	widgets.QApplication

	Search *SearchUI
	Login  *LoginUI
}

func setupUI() {
	// NewGui will create a new QApplication
	ui = NewUI(len(os.Args), os.Args)

	// setup the windows
	ui.Search = setupSearch()
	ui.Login = setupLogin()
}
