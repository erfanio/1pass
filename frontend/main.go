package frontend

import (
	"github.com/therecipe/qt/widgets"
	"os"
)

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

var app *widgets.QApplication

func InitGui() {
	// new app
	app = widgets.NewQApplication(len(os.Args), os.Args)
	app.SetStyleSheet(stylesheet)
}

func StartGui() {
	app.Exec()
}

func CloseGui() {
	app.Quit()
}

func ShowError(msg string) {
	errorDialog := widgets.NewQErrorMessage(nil)
	errorDialog.ShowMessage(msg)
	errorDialog.Exec()
}
