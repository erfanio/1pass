package main

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	settings *core.QSettings
)

func main() {
	// setup
	settings = core.NewQSettings2(core.QSettings__UserScope, "erfan.io", "1pass", nil)

	setupUI()
	ui.Search.Show()

	tryFetching()

	// start app
	ui.Exec()
}

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

	// get list of item summaries from 1pass's cli
	cmd := exec.Command("/bin/op", "list", "items")
	cmd.Env = append(os.Environ(), session)
	output, err := cmd.Output()

	// login if error (probably not logged in)
	if err != nil {
		log.Print(err)
		promptLogin()
	} else {
		items := json2list(output)
		populateList(items)
	}
}
