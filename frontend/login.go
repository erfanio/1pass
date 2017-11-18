package frontend

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Login struct {
	widgets.QDialog

	_ func() `constructor:"init"`

	layout   *widgets.QFormLayout
	domain   *widgets.QLineEdit
	email    *widgets.QLineEdit
	key      *widgets.QLineEdit
	password *widgets.QLineEdit
	button   *widgets.QPushButton
}

func (w *Login) init() {
	// login is triggered by clicking enter on inputs
	// window needs to be a dialog for this work
	w.SetWindowTitle("Login")

	// vertical layout
	w.layout = widgets.NewQFormLayout(nil)
	w.SetLayout(w.layout)

	// add the form (inputs and a button)
	w.domain = widgets.NewQLineEdit(nil)
	w.domain.SetText("my.1password.com")
	w.layout.AddRow3("Domain", w.domain)

	w.email = widgets.NewQLineEdit(nil)
	w.layout.AddRow3("Email", w.email)

	w.key = widgets.NewQLineEdit(nil)
	w.layout.AddRow3("Secret Key", w.key)

	w.password = widgets.NewQLineEdit(nil)
	// don't show the password
	w.password.SetEchoMode(widgets.QLineEdit__Password)
	w.layout.AddRow3("Master Password", w.password)

	w.button = widgets.NewQPushButton2("Sign In", nil)
	w.button.SetDefault(true)
	w.layout.AddRow5(w.button)

	// focus on first empty input
	w.ConnectShowEvent(func(event *gui.QShowEvent) {
		// a list of inputs in order
		inputs := []*widgets.QLineEdit{w.domain, w.email, w.key, w.password}

		for i := range inputs {
			input := inputs[i]
			if input.Text() == "" {
				input.SetFocus(core.Qt__NoFocusReason)
				return
			}
		}
	})

	// close the app if login is rejected (esc key)
	w.ConnectRejected(func() {
		CloseGui()
	})
}

func (w *Login) OnSubmit(listener func(string, string, string, string)) {
	w.button.ConnectClicked(func(checked bool) {
		domainText := w.domain.Text()
		emailText := w.email.Text()
		keyText := w.key.Text()
		passwordText := w.password.Text()

		listener(domainText, emailText, keyText, passwordText)
	})
}

func (w *Login) SetDomain(text string) {
	w.domain.SetText(text)
}

func (w *Login) SetEmail(text string) {
	w.email.SetText(text)
}

func (w *Login) SetKey(text string) {
	w.key.SetText(text)
}

func (w *Login) SetPassword(text string) {
	w.password.SetText(text)
}

func (w *Login) StartWait() {
	cursor := gui.NewQCursor2(core.Qt__WaitCursor)
	w.SetCursor(cursor)
	w.button.SetDisabled(true)
	// this is necessary to process this event instantly
	app.ProcessEvents(core.QEventLoop__AllEvents)
}

func (w *Login) EndWait() {
	cursor := gui.NewQCursor2(core.Qt__ArrowCursor)
	w.SetCursor(cursor)
	w.button.SetDisabled(false)
	// this is necessary to process this event instantly
	app.ProcessEvents(core.QEventLoop__AllEvents)
}
