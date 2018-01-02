package main

import (
	"fmt"
	"github.com/erfanio/1pass/ui"
	"github.com/therecipe/qt/core"
	"log"
	"os"
	"os/exec"
	"strings"
)

func getSession() string {
	// get logged in session env
	domain := settings.Value("domain", core.NewQVariant17("my.1password.com")).ToString()
	subdomain := strings.Split(domain, ".")[0]
	sessionKey := settings.Value("session_key", core.NewQVariant17("")).ToString()
	// e.g. OP_SESSION_my=abcdefg
	return fmt.Sprintf("OP_SESSION_%v=%v", subdomain, sessionKey)
}

func tryFetching() {
	// if logged in, session will let us send queries for 30 min
	session := getSession()

	go func() {
		// get list of item summaries from 1pass's cli
		cmd := exec.Command("/bin/op", "list", "items")
		cmd.Env = append(os.Environ(), session)
		output, err := cmd.Output()

		// login if error (probably not logged in)
		if err != nil {
			logOpError(err)
			promptLogin()
		} else {
			items := json2list(output)
			populateList(items)
		}
	}()
}

func fetchDetails(UUID string, callback func(completeItem)) {
	ui.App.Search.Disable()
	// if logged in, session will let us send queries for 30 min
	session := getSession()

	go func() {
		// get list of item summaries from 1pass's cli
		cmd := exec.Command("/bin/op", "get", "item", UUID)
		cmd.Env = append(os.Environ(), session)
		output, err := cmd.Output()
		ui.App.Search.EnableAndFocus()

		// login if error (probably not logged in)
		if err != nil {
			logOpError(err)
			promptLogin()
		} else {
			item := json2item(output)
			callback(item)
		}
	}()
}

func submitLogin(domain, email, key, password string) {
	ui.App.Login.Disable()

	go func() {
		// try to login and get the session token (raw outputs only the token)
		out, err := exec.Command("/bin/op", "signin", domain, email, key, password, "--output=raw").Output()
		ui.App.Login.Enable()

		if err != nil {
			logOpError(err)

			msg := "Login failed!\n"
			// try get the stderr to show in dialog (if it's an ExitError)
			exitErr, ok := err.(*exec.ExitError)
			if ok {
				msg += string(exitErr.Stderr)
			} else {
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

func logOpError(err error) {
	// try getting the strderr returned
	exitErr, ok := err.(*exec.ExitError)
	if ok {
		log.Println(err, string(exitErr.Stderr))
	} else {
		log.Println(err)
	}
}
