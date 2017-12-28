package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const (
	loginWidth        = 400
	loginInputHeight  = 35
	loginButtonWidth  = 180
	loginButtonHeight = 40
)

const loginStyles = `
QLineEdit, QPushButton {
  font-size: 12px;
  background-color: #FFFFFF;
  border-width: 1px;
  border-style: solid;
  border-color: #BDBDBD;
	padding: 5px;
}

QLineEdit {
  border-radius: 2px;
}

QPushButton {
  border-radius: 10px;
}

QPushButton:hover {
	background-color: #F6F7FB;
}
`

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

	Layout      *widgets.QVBoxLayout
	Domain      *widgets.QLineEdit
	Email       *widgets.QLineEdit
	Key         *widgets.QLineEdit
	Password    *widgets.QLineEdit
	Button      *widgets.QPushButton
	ErrorDialog *widgets.QErrorMessage
}

func setupLogin() *LoginUI {
	// create the window
	w := NewLoginUI(nil, core.Qt__Tool|core.Qt__WindowStaysOnTopHint)
	w.SetWindowTitle("Login")
	w.SetStyleSheet(loginStyles)

	// vertical layout
	w.Layout = widgets.NewQVBoxLayout()
	w.SetLayout(w.Layout)

	// add the form (inputs and a button)
	w.Layout.AddWidget(widgets.NewQLabel2("Domain", nil, core.Qt__Widget), 0, 0)
	w.Domain = widgets.NewQLineEdit(nil)
	w.Domain.SetText("my.1password.com")
	w.Domain.SetFixedHeight(loginInputHeight)
	w.Layout.AddWidget(w.Domain, 0, 0)

	w.Layout.AddWidget(widgets.NewQLabel2("Email", nil, core.Qt__Widget), 0, 0)
	w.Email = widgets.NewQLineEdit(nil)
	w.Email.SetFixedHeight(loginInputHeight)
	w.Layout.AddWidget(w.Email, 0, 0)

	w.Layout.AddWidget(widgets.NewQLabel2("Secret Key", nil, core.Qt__Widget), 0, 0)
	w.Key = widgets.NewQLineEdit(nil)
	w.Key.SetFixedHeight(loginInputHeight)
	w.Layout.AddWidget(w.Key, 0, 0)

	w.Layout.AddWidget(widgets.NewQLabel2("Master Password", nil, core.Qt__Widget), 0, 0)
	w.Password = widgets.NewQLineEdit(nil)
	w.Password.SetFixedHeight(loginInputHeight)
	// don't show the password
	w.Password.SetEchoMode(widgets.QLineEdit__Password)
	w.Layout.AddWidget(w.Password, 0, 0)

	w.Button = widgets.NewQPushButton2("Sign In", nil)
	w.Button.SetDefault(true)
	w.Button.SetSizePolicy2(widgets.QSizePolicy__Maximum, widgets.QSizePolicy__Preferred)
	w.Layout.AddWidget(w.Button, 0, 0)

	w.ErrorDialog = widgets.NewQErrorMessage(nil)

	w.setupEventListeners()
	return w
}

func (w *LoginUI) setupEventListeners() {
	// set width hint
	w.ConnectSizeHint(func() *core.QSize {
		return core.NewQSize2(loginWidth, w.SizeHintDefault().Height())
	})

	w.Button.ConnectSizeHint(func() *core.QSize {
		return core.NewQSize2(loginButtonWidth, loginButtonHeight)
	})

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
		println(core.QThread_CurrentThread().Pointer())
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

// Enable is a shortcut for SetDisabled(false) and setting cursor back to default
func (w *LoginUI) Enable() {
	w.SetDisabled(false)
	cursor := gui.NewQCursor2(core.Qt__ArrowCursor)
	w.SetCursor(cursor)
}

// Disable is a shortcut for SetDisabled(true) and setting cursor to wait
func (w *LoginUI) Disable() {
	w.SetDisabled(true)
	cursor := gui.NewQCursor2(core.Qt__WaitCursor)
	w.SetCursor(cursor)
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
	w.ErrorDialog.ShowMessage(msg)
	w.ErrorDialog.Exec()
}
