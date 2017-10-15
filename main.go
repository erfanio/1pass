package main

import (
	"fmt"
)

func main() {
	main := initGui()
	main.window.Show()

	login := initLogin(func(key, password string) {
		fmt.Printf("Key: %v\nPassword: %v\n", key, password)
	})
	login.window.Show()

	// start the app
	main.app.Exec()
}
