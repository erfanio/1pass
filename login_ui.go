package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// LoginUI is a "class" that is the login dialog
// it stores the Qt objects and provides the slots to manipulate it
type LoginUI struct {
	widgets.QDialog

	_ func()                               `slots:"Show"`
	_ func()                               `slots:"Hide"`
	_ func()                               `slots:"Disable"`
	_ func()                               `slots:"Enable"`
	_ func(bool)                           `slots:"SetDisabled"`
	_ func(string, string, string, string) `slots:"SetInputTexts"`
	_ func(string)                         `slots:"ShowError"`

	Layout   *widgets.QFormLayout
	Domain   *widgets.QLineEdit
	Email    *widgets.QLineEdit
	Key      *widgets.QLineEdit
	Password *widgets.QLineEdit
	Button   *widgets.QPushButton
}

func setupLogin() *LoginUI {
	// create the window
	w := NewLoginUI(nil, core.Qt__Tool|core.Qt__WindowStaysOnTopHint)
	w.SetWindowTitle("Login")

	// vertical layout
	w.Layout = widgets.NewQFormLayout(nil)
	w.SetLayout(w.Layout)

	// add the form (inputs and a button)
	w.Domain = widgets.NewQLineEdit(nil)
	w.Domain.SetText("my.1password.com")
	w.Layout.AddRow3("Domain", w.Domain)

	w.Email = widgets.NewQLineEdit(nil)
	w.Layout.AddRow3("Email", w.Email)

	w.Key = widgets.NewQLineEdit(nil)
	w.Layout.AddRow3("Secret Key", w.Key)

	w.Password = widgets.NewQLineEdit(nil)
	// don't show the password
	w.Password.SetEchoMode(widgets.QLineEdit__Password)
	w.Layout.AddRow3("Master Password", w.Password)

	w.Button = widgets.NewQPushButton2("Sign In", nil)
	w.Button.SetDefault(true)
	w.Layout.AddRow5(w.Button)

	w.setupEventListeners()
	return w
}

func (w *LoginUI) setupEventListeners() {
	// focus on first empty input
	w.ConnectShowEvent(func(event *gui.QShowEvent) {
		// a list of inputs in order
		inputs := []*widgets.QLineEdit{w.Domain, w.Email, w.Key, w.Password}

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
		ui.Quit()
	})

	// listen for form submission
	w.Button.ConnectClicked(func(checked bool) {
		submitLogin(w.Domain.Text(), w.Email.Text(), w.Key.Text(), w.Password.Text())
	})
}

// Show shows the login dialog
func (w *LoginUI) Show() {
	w.QDialog.Show()
}

// Hide hides the login dialog
func (w *LoginUI) Hide() {
	w.QDialog.Hide()
}

// SetDisabled sets the disabled state of the dialog components (inputs, button)
func (w *LoginUI) SetDisabled(disabled bool) {
	w.Domain.SetDisabled(disabled)
	w.Email.SetDisabled(disabled)
	w.Key.SetDisabled(disabled)
	w.Password.SetDisabled(disabled)
	w.Button.SetDisabled(disabled)
}

// Enable is a shortcut for SetDisabled(false)
func (w *LoginUI) Enable() {
	w.SetDisabled(false)
}

// Disable is a shortcut for SetDisabled(true)
func (w *LoginUI) Disable() {
	w.SetDisabled(true)
}

// SetInputTexts sets the text in the inputs (if string is empty will not override previous state)
func (w *LoginUI) SetInputTexts(domain, email, key, password string) {
	if len(domain) > 0 {
		w.Domain.SetText(domain)
	}
	if len(email) > 0 {
		w.Email.SetText(email)
	}
	if len(key) > 0 {
		w.Key.SetText(key)
	}
	if len(password) > 0 {
		w.Password.SetText(password)
	}
}

// ShowError will display a dismissable error dialog with a message
func (w *LoginUI) ShowError(msg string) {
	errorDialog := widgets.NewQErrorMessage(nil)
	errorDialog.ShowMessage(msg)
	errorDialog.Exec()
}
