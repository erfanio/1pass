package main

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"os"
	"strconv"
)

const (
	APP_NAME        = "1Password Lookup"
	WINDOW_WIDTH    = 600
	EDITLINE_HEIGHT = 50
	RESULT_HEIGHT   = 50
)
const stylesheet = `
* {
  background-color: #EEEEEE;
}

#innerWindow {
  padding: 5px;
  border-radius: 5px;
  border-width: 1px;
  border-style: solid;
  border-color: #E0E0E0;
}

#input {
  font-size: 28px;
  background-color: #fff;
  border-radius: 2px;
  border-width: 1px;
  border-style: solid;
  border-color: #E0E0E0;
}

#list {
  border: none;
  selection-background-color: #40C4FF;
  selection-color: #000;
}

`

func main() {
	fmt.Println("hi")
	// new app
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetStyleSheet(stylesheet)

	// main window is floating
	window := widgets.NewQWidget(nil, core.Qt__Tool|
		core.Qt__FramelessWindowHint|
		core.Qt__WindowCloseButtonHint|
		core.Qt__WindowStaysOnTopHint)
	window.SetWindowTitle(APP_NAME)
	// tell window to quit when it closes (Qt::Tool turns this off for some reason)
	window.SetAttribute(core.Qt__WA_QuitOnClose, true)
	window.SetAttribute(core.Qt__WA_TranslucentBackground, true)

	// window is layed out vertically
	windowLayout := widgets.NewQVBoxLayout()
	windowLayout.SetContentsMargins(0, 0, 0, 0)
	window.SetLayout(windowLayout)

	// add a inner window widget (since window is completely transparent this is needed for the border)
	innerWindow := widgets.NewQFrame(nil, 0)
	innerWindow.SetObjectName("innerWindow")
	windowLayout.AddWidget(innerWindow, 0, 0)

	// inner window is layed vertically
	layout := widgets.NewQVBoxLayout()
	layout.SetContentsMargins(0, 0, 0, 0)
	innerWindow.SetLayout(layout)

	// LineEdit for the search query
	input := widgets.NewQLineEdit(nil)
	input.SetObjectName("input")
	input.SetFixedHeight(EDITLINE_HEIGHT)
	layout.AddWidget(input, 0, 0)

	// items in the list of logins
	item := widgets.NewQStyledItemDelegate(nil)

	// list of logins
	list := widgets.NewQListWidget(nil)
	list.SetObjectName("list")
	// expands horizontally but sticks to size hint vertically
	list.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Fixed)
	// set the list items
	list.SetItemDelegate(item)
	// hidden by default (shown as soon as something appears inside)
	list.Hide()
	layout.AddWidget(list, 0, 0)

	// exit on esc
	input.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		if event.Key() == int(core.Qt__Key_Escape) {
			app.Quit()
		} else {
			input.KeyPressEventDefault(event)
		}
	})

	// handle text change
	input.ConnectTextChanged(func(text string) {
		items := make([]string, len(text))
		// fake items for now
		for i := 0; i < len(text); i++ {
			items = append(items, strconv.Itoa(i)+" "+text[:i+1])
		}

		// replace the list with current items (or hide it if empty)
		list.Clear()
		if len(items) > 0 {
			list.Show()
			list.AddItems(items)
			list.UpdateGeometry()
		} else {
			list.Hide()
		}

		// update the size of the parents
		parent := list.ParentWidget()
		for parent != nil {
			parent.AdjustSize()
			parent = parent.ParentWidget()
		}
	})

	// size of the list items
	item.ConnectSizeHint(func(option *widgets.QStyleOptionViewItem, index *core.QModelIndex) *core.QSize {
		// only care about setting the correct height
		return core.NewQSize2(item.SizeHintDefault(option, index).Width(), RESULT_HEIGHT)
	})
	// // draw the list item
	// resultItem.ConnectPaint(func(painter *gui.QPainter, option *widgets.QStyleOptionViewItem, index *core.QModelIndex) {
	//   data := index.Data(int(core.Qt__DisplayRole))
	//   fmt.Println(data.ToString(), index.Row())
	// })
	// let the list be empty (shrink to 0px)
	list.ConnectMinimumSizeHint(func() *core.QSize {
		return core.NewQSize2(0, 0)
	})
	// size of the list (don't let the list grow bigger than 5 items)
	list.ConnectSizeHint(func() *core.QSize {
		rowSize := list.SizeHintForRow(0)
		// -1 when hidden should be 0
		if rowSize < 0 {
			rowSize = 0
		}
		count := list.Count()
		if count > 5 {
			count = 5
		}
		// only care about setting the correct height
		return core.NewQSize2(list.SizeHintDefault().Width(), RESULT_HEIGHT*count)
	})

	// run the app
	window.Show()
	app.Exec()
}
