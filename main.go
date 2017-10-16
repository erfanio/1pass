package main

import (
	"fmt"
	"github.com/erfanio/1pass/frontend"
	"github.com/therecipe/qt/core"
	"os/exec"
)

func main() {
	// app
	frontend.InitGui()
	settings := core.NewQSettings2(core.QSettings__UserScope, "erfan.io", "1pass", nil)

	// search window
	search := frontend.NewSearch()
	search.Show()

	// login
	login := frontend.NewLogin()
	// get input's previous state (or default)
	login.SetDomain(settings.Value("domain", core.NewQVariant17("my.1password.com")).ToString())
	login.SetEmail(settings.Value("email", core.NewQVariant17("")).ToString())
	login.SetKey(settings.Value("key", core.NewQVariant17("")).ToString())

	// listener for login
	login.OnSubmit(func(domain, email, key, password string) {
		// remember the state for next time
		settings.SetValue("domain", core.NewQVariant17(domain))
		settings.SetValue("email", core.NewQVariant17(email))
		settings.SetValue("key", core.NewQVariant17(key))

		// put into loading state
		login.StartWait()

		// try to login and get the session token (raw outputs only the token)
		out, err := exec.Command("/bin/op", "signin", domain, email, key, password, "--output=raw").Output()

		// end loading state
		login.EndWait()
		if err != nil {
			// show error in a dialog
			msg := "Login failed! "
			exitErr, ok := err.(*exec.ExitError)
			if ok {
				msg += string(exitErr.Stderr)
			}
			frontend.ShowError(msg)
			return
		}

		// successful login so hide login window
		login.Hide()

		fmt.Println(string(out))
	})
	login.Show()

	// start app
	frontend.StartGui()
}
