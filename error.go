package main

import (
	"github.com/therecipe/qt/widgets"
)

func showError(msg string) {
	errorDialog := widgets.NewQErrorMessage(nil)
	errorDialog.ShowMessage(msg)
	errorDialog.Exec()
}
