package ui

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
QLineEdit, QPushButton, QLabel {
  font-size: 12px;
}

QLineEdit, QPushButton {
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

	_ func()                               `constructor:"init"`
	_ func()                               `slot:"Start"` // start/finish because show/hide would collide with QDialog's show/hide
	_ func()                               `slot:"Finish"`
	_ func()                               `slot:"Disable"`
	_ func()                               `slot:"Enable"`
	_ func(string, string, string, string) `slot:"SetInputTexts"`
	_ func(string)                         `slot:"ShowError"`

	layout      *widgets.QVBoxLayout
	domain      *widgets.QLineEdit
	email       *widgets.QLineEdit
	key         *widgets.QLineEdit
	password    *widgets.QLineEdit
	button      *widgets.QPushButton
	errorDialog *widgets.QErrorMessage

	loginListener func(string, string, string, string)
}

func (w *LoginUI) init() {
	// create the window
	w.SetWindowTitle("Login")
	w.SetStyleSheet(loginStyles)

	// vertical layout
	w.layout = widgets.NewQVBoxLayout()
	w.SetLayout(w.layout)

	// add the form (inputs and a button)
	w.layout.AddWidget(widgets.NewQLabel2("Domain", nil, core.Qt__Widget), 0, 0)
	w.domain = widgets.NewQLineEdit(nil)
	w.domain.SetText("my.1password.com")
	w.domain.SetFixedHeight(loginInputHeight)
	w.layout.AddWidget(w.domain, 0, 0)

	w.layout.AddWidget(widgets.NewQLabel2("Email", nil, core.Qt__Widget), 0, 0)
	w.email = widgets.NewQLineEdit(nil)
	w.email.SetFixedHeight(loginInputHeight)
	w.layout.AddWidget(w.email, 0, 0)

	w.layout.AddWidget(widgets.NewQLabel2("Secret Key", nil, core.Qt__Widget), 0, 0)
	w.key = widgets.NewQLineEdit(nil)
	w.key.SetFixedHeight(loginInputHeight)
	w.key.SetEchoMode(widgets.QLineEdit__Password)
	w.layout.AddWidget(w.key, 0, 0)

	w.layout.AddWidget(widgets.NewQLabel2("Master Password", nil, core.Qt__Widget), 0, 0)
	w.password = widgets.NewQLineEdit(nil)
	w.password.SetFixedHeight(loginInputHeight)
	w.password.SetEchoMode(widgets.QLineEdit__Password)
	w.layout.AddWidget(w.password, 0, 0)

	w.button = widgets.NewQPushButton2("Sign In", nil)
	w.button.SetDefault(true)
	w.button.SetSizePolicy2(widgets.QSizePolicy__Maximum, widgets.QSizePolicy__Preferred)
	w.layout.AddWidget(w.button, 0, 0)

	w.errorDialog = widgets.NewQErrorMessage(nil)

	// set sizes
	w.ConnectSizeHint(func() *core.QSize {
		return core.NewQSize2(loginWidth, w.SizeHintDefault().Height())
	})
	w.button.ConnectSizeHint(func() *core.QSize {
		return core.NewQSize2(loginButtonWidth, loginButtonHeight)
	})

	w.ConnectShowEvent(w.focusInput)
	// close the app if login is rejected (esc key)
	w.ConnectRejected(func() {
		App.Quit()
	})
	// listen for form submission
	w.button.ConnectClicked(func(checked bool) {
		if w.loginListener != nil {
			w.loginListener(w.domain.Text(), w.email.Text(), w.key.Text(), w.password.Text())
		}
	})

	// Setup slots
	w.ConnectStart(func() {
		w.Show()
	})
	w.ConnectFinish(func() {
		w.Hide()
	})
	w.ConnectEnable(w.enable)
	w.ConnectDisable(w.disable)
	w.ConnectSetInputTexts(w.setInputTexts)
	w.ConnectShowError(w.showError)
}

func (w *LoginUI) SetLoginListener(f func(string, string, string, string)) {
	w.loginListener = f
}

// focusInput focuses on the first empty input
func (w *LoginUI) focusInput(event *gui.QShowEvent) {
	// a list of inputs in order
	inputs := []*widgets.QLineEdit{w.domain, w.email, w.key, w.password}

	for i := range inputs {
		input := inputs[i]
		if input.Text() == "" {
			input.SetFocus(core.Qt__NoFocusReason)
			return
		}
	}
}

// setDisabled sets the disabled state of the dialog components (inputs, button)
func (w *LoginUI) setDisabled(disabled bool) {
	w.domain.SetDisabled(disabled)
	w.email.SetDisabled(disabled)
	w.key.SetDisabled(disabled)
	w.password.SetDisabled(disabled)
	w.button.SetDisabled(disabled)
}

// enable is a shortcut for SetDisabled(false) and setting cursor back to default
func (w *LoginUI) enable() {
	w.setDisabled(false)
	cursor := gui.NewQCursor2(core.Qt__ArrowCursor)
	w.SetCursor(cursor)
}

// disable is a shortcut for SetDisabled(true) and setting cursor to wait
func (w *LoginUI) disable() {
	w.setDisabled(true)
	cursor := gui.NewQCursor2(core.Qt__WaitCursor)
	w.SetCursor(cursor)
}

// setInputTexts sets the text in the inputs (if string is empty will not override previous state)
func (w *LoginUI) setInputTexts(domain, email, key, password string) {
	if len(domain) > 0 {
		w.domain.SetText(domain)
	}
	if len(email) > 0 {
		w.email.SetText(email)
	}
	if len(key) > 0 {
		w.key.SetText(key)
	}
	if len(password) > 0 {
		w.password.SetText(password)
	}
}

// showError will display a dismissable error dialog with a message
func (w *LoginUI) showError(msg string) {
	w.errorDialog.ShowMessage(msg)
	w.errorDialog.Exec()
}
