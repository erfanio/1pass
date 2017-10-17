package main

import (
	"fmt"
	"github.com/erfanio/1pass/frontend"
	"github.com/therecipe/qt/core"
	"os"
	"os/exec"
	"strings"
)

var (
	settings *core.QSettings
	search   frontend.Search
	login    frontend.Login
)

func main() {
	// app
	frontend.InitGui()
	settings = core.NewQSettings2(core.QSettings__UserScope, "erfan.io", "1pass", nil)

	// search window
	search = frontend.NewSearch()
	search.Show()

	// get the list of items (if fails, will prompt login)
	getList()

	// start app
	frontend.StartGui()
}

func getList() {
	// get items from op using the session key
	domain := settings.Value("domain", core.NewQVariant17("my.1password.com")).ToString()
	session_key := settings.Value("session_key", core.NewQVariant17("")).ToString()

	subdomain := strings.Split(domain, ".")[0]
	// e.g. OP_SESSION_my=abcdefg
	session_env := fmt.Sprintf("OP_SESSION_%v=%v", subdomain, session_key)

	cmd := exec.Command("/bin/op", "list", "items")
	cmd.Env = append(os.Environ(), session_env)
	out, err := cmd.Output()

	// login if list couldn't be fetched
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			fmt.Println(err.Error())
			fmt.Println(string(exitErr.Stderr))
		}
		initLogin()
		return
	}

	fmt.Println(string(out))
}

func initLogin() {
	// login
	login = frontend.NewLogin()
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

		// remember the session key
		session := strings.TrimSpace(string(out))
		settings.SetValue("session_key", core.NewQVariant17(session))
	})

	login.Show()
}
