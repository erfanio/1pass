package main

import (
	"github.com/therecipe/qt/core"
	"log"
	"os/exec"
	"strings"
)

// show the login form
func promptLogin() {
	// populate form with previously stored info
	ui.Login.SetInputTexts(
		settings.Value("domain", core.NewQVariant17("my.1password.com")).ToString(),
		settings.Value("email", core.NewQVariant17("")).ToString(),
		settings.Value("key", core.NewQVariant17("")).ToString(),
		"")

	ui.Login.Show()
}

func submitLogin(domain, email, key, password string) {
	// try to login and get the session token (raw outputs only the token)
	out, err := exec.Command("/bin/op", "signin", domain, email, key, password, "--output=raw").Output()

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
		showError(msg)
	} else {
		ui.Login.Hide()

		// save session (valid for 30 min)
		session := strings.TrimSpace(string(out))
		settings.SetValue("session_key", core.NewQVariant17(session))
	}
}
