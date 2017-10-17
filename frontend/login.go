package frontend

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Login struct {
	window   *widgets.QDialog
	layout   *widgets.QFormLayout
	domain   *widgets.QLineEdit
	email    *widgets.QLineEdit
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
	layout := widgets.NewQFormLayout(nil)
	window.SetLayout(layout)

	// add the form (inputs and a button)
	domain := widgets.NewQLineEdit(nil)
	domain.SetText("my.1password.com")
	layout.AddRow3("Domain", domain)

	email := widgets.NewQLineEdit(nil)
	layout.AddRow3("Email", email)

	key := widgets.NewQLineEdit(nil)
	layout.AddRow3("Secret Key", key)

	password := widgets.NewQLineEdit(nil)
	// don't show the password
	password.SetEchoMode(widgets.QLineEdit__Password)
	layout.AddRow3("Master Password", password)

	button := widgets.NewQPushButton2("Sign In", nil)
	button.SetDefault(true)
	layout.AddRow5(button)

	// focus on first empty input
	window.ConnectShowEvent(func(event *gui.QShowEvent) {
		// a list of inputs in order
		inputs := []*widgets.QLineEdit{domain, email, key, password}

		for i := range inputs {
			input := inputs[i]
			if input.Text() == "" {
				input.SetFocus(core.Qt__NoFocusReason)
				return
			}
		}
	})

	// close the app if login is rejected (esc key)
	window.ConnectRejected(func() {
		CloseGui()
	})

	// callback on click
	return Login{
		window,
		layout,
		domain,
		email,
		key,
		password,
		button,
	}
}

func (l *Login) OnSubmit(listener func(string, string, string, string)) {
	l.button.ConnectClicked(func(checked bool) {
		domainText := l.domain.Text()
		emailText := l.email.Text()
		keyText := l.key.Text()
		passwordText := l.password.Text()

		listener(domainText, emailText, keyText, passwordText)
	})
}

func (l *Login) Show() {
	l.window.Show()
}

func (l *Login) Hide() {
	l.window.Hide()
}

func (l *Login) SetDomain(text string) {
	l.domain.SetText(text)
}

func (l *Login) SetEmail(text string) {
	l.email.SetText(text)
}

func (l *Login) SetKey(text string) {
	l.key.SetText(text)
}

func (l *Login) SetPassword(text string) {
	l.password.SetText(text)
}

func (l *Login) StartWait() {
	cursor := gui.NewQCursor2(core.Qt__WaitCursor)
	l.window.SetCursor(cursor)
	l.button.SetDisabled(true)
	// this is necessary to process this event instantly
	app.ProcessEvents(core.QEventLoop__AllEvents)
}

func (l *Login) EndWait() {
	cursor := gui.NewQCursor2(core.Qt__ArrowCursor)
	l.window.SetCursor(cursor)
	l.button.SetDisabled(false)
	// this is necessary to process this event instantly
	app.ProcessEvents(core.QEventLoop__AllEvents)
}
