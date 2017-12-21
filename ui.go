package main

import (
	"github.com/therecipe/qt/widgets"
	"os"
)

var ui *UI

type UI struct {
	widgets.QApplication

	Search *SearchUI
	Login  *LoginUI
}

const stylesheet = `
* {
  background-color: #EEEEEE;
}

#innerWindow {
  padding: 5px;
  border-radius: 5px;
  border-width: 1px;
  border-style: solid;
  border-color: #E0E0E0;
}

#input {
  font-size: 28px;
  background-color: #fff;
  border-radius: 2px;
  border-width: 1px;
  border-style: solid;
  border-color: #E0E0E0;
}

#list {
  border: none;
  selection-background-color: #40C4FF;
  selection-color: #000;
}
`

func setupUI() {
	// NewGui will create a new QApplication
	ui = NewUI(len(os.Args), os.Args)
	ui.SetStyleSheet(stylesheet)

	// setup the windows
	ui.Search = setupSearch()
	ui.Login = setupLogin()
}
