package main

import (
	"fmt"
	"github.com/erfanio/1pass/frontend"
	"os/exec"
)

func main() {
	// app
	frontend.InitGui()

	// search window
	search := frontend.NewSearch()
	search.Show()

	// login
	login := frontend.NewLogin()
	// listener for login
	login.OnSubmit(func(domain, email, key, password string) {
		// put into loading state
		login.StartWait()

		// try to login and get the session token (raw outputs only the token)
		command := fmt.Sprintf("op signin %v %v %v %v --output=raw", domain, email, key, password)
		out, err := exec.Command("/bin/bash", "-c", command).Output()

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
