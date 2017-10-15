package main

func main() {
	main := initGui()
	main.window.Show()

	login := initLogin()
	login.window.Show()

	// start the app
	main.app.Exec()
}
