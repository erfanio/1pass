package main

import (
	"github.com/erfanio/1pass/ui"
	"github.com/therecipe/qt/core"
	"log"
	"os/exec"
	"strings"
)

func setupLogin() {
	ui.App.Login.SetLoginListener(submitLogin)
}

// show the login form
func promptLogin() {
	// populate form with previously stored info
	ui.App.Login.SetInputTexts(
		settings.Value("domain", core.NewQVariant17("my.1password.com")).ToString(),
		settings.Value("email", core.NewQVariant17("")).ToString(),
		settings.Value("key", core.NewQVariant17("")).ToString(),
		"")

	ui.App.Login.Start()
}

func submitLogin(domain, email, key, password string) {
	ui.App.Login.Disable()

	go func() {
		// try to login and get the session token (raw outputs only the token)
		out, err := exec.Command("/bin/op", "signin", domain, email, key, password, "--output=raw").Output()
		ui.App.Login.Enable()

		if err != nil {
			msg := "Login failed!\n"
			// try get the stderr to show in dialog (if it's an ExitError)
			exitErr, ok := err.(*exec.ExitError)
			if ok {
				log.Print(err, string(exitErr.Stderr))
				msg += string(exitErr.Stderr)
			} else {
				log.Print(err)
				msg += "Unexpected error... Please create an issue on github :)"
			}

			// show error (and the stderr) in a dialog
			ui.App.Login.ShowError(msg)
		} else {
			ui.App.Login.Finish()

			// save session (valid for 30 min)
			session := strings.TrimSpace(string(out))
			settings.SetValue("session_key", core.NewQVariant17(session))

			tryFetching()
		}
	}()
}
