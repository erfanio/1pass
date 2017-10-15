package main

import (
	"fmt"
	"github.com/erfanio/1pass/frontend"
)

func main() {
	frontend.InitGui()
	search := frontend.NewSearch()
	search.Show()

	login := frontend.NewLogin()
	login.OnSubmit(func(key, password string) {
		fmt.Printf("Key: %v\nPassword: %v\n", key, password)
		login.Hide()
	})
	login.Show()

	// start the app
	frontend.StartGui()
}
