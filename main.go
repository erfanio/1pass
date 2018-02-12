package main

import (
	"fmt"
	"github.com/erfanio/1pass/ui"
	"github.com/therecipe/qt/core"
)

var (
	settings *core.QSettings
)

func main() {
	fmt.Println("hi :)")
	// setup
	settings = core.NewQSettings2(core.QSettings__UserScope, "erfan.io", "1pass", nil)

	ui.SetupUI()
	ui.App.Search.Start()
	setupList()
	ui.App.Login.SetLoginListener(submitLogin)

	tryFetching()

	// start app
	ui.App.Exec()
}

// show the login form
func promptLogin() {
	// populate form with previously stored info
	ui.App.Login.SetInputTexts(
		settings.Value("domain", core.NewQVariant17("my.1password.com")).ToString(),
		settings.Value("email", core.NewQVariant17("")).ToString(),
		settings.Value("key", core.NewQVariant17("")).ToString(),
		"",
	)
	ui.App.Login.Start()
}
