package frontend

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type Login struct {
	window   *widgets.QDialog
	layout   *widgets.QVBoxLayout
	key      *widgets.QLineEdit
	password *widgets.QLineEdit
	button   *widgets.QPushButton
}

func NewLogin() Login {
	// login is triggered by clicking enter on inputs
	// window needs to be a dialog for this work
	window := widgets.NewQDialog(nil, core.Qt__Tool|core.Qt__WindowStaysOnTopHint)
	window.SetWindowTitle("Login")

	// vertical layout
	layout := widgets.NewQVBoxLayout()
	window.SetLayout(layout)

	// add the form (inputs and a button)
	key := widgets.NewQLineEdit(nil)
	key.SetPlaceholderText("Secret Key")
	layout.AddWidget(key, 0, 0)

	password := widgets.NewQLineEdit(nil)
	password.SetPlaceholderText("Password")
	// don't show the password
	password.SetEchoMode(widgets.QLineEdit__Password)
	layout.AddWidget(password, 0, 0)

	button := widgets.NewQPushButton2("Login", nil)
	button.SetDefault(true)
	layout.AddWidget(button, 0, 0)

	// callback on click
	return Login{
		window,
		layout,
		key,
		password,
		button,
	}
}

func (l *Login) OnSubmit(onSubmit func(string, string)) {
	l.button.ConnectClicked(func(checked bool) {
		keyText := l.key.Text()
		passwordText := l.password.Text()

		onSubmit(keyText, passwordText)
	})
}

func (l *Login) Show() {
  l.window.Show()
}

func (l *Login) Hide() {
  l.window.Hide()
}
