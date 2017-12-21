package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const (
	APP_NAME        = "1Password Lookup"
	WINDOW_WIDTH    = 600
	EDITLINE_HEIGHT = 50
	RESULT_HEIGHT   = 50
)

type SearchUI struct {
	widgets.QWidget

	WindowLayout *widgets.QVBoxLayout
	InnerWindow  *widgets.QFrame
	Layout       *widgets.QVBoxLayout
	Input        *widgets.QLineEdit
	Item         *widgets.QStyledItemDelegate
	List         *widgets.QListView
}

func setupSearch() *SearchUI {
	// create the window
	w := NewSearchUI(nil, core.Qt__Tool|
		core.Qt__FramelessWindowHint|
		core.Qt__WindowCloseButtonHint|
		core.Qt__WindowStaysOnTopHint)
	w.SetWindowTitle(APP_NAME)

	// tell window to quit when it closes (Qt::Tool turns this off for some reason)
	w.SetAttribute(core.Qt__WA_QuitOnClose, true)
	w.SetAttribute(core.Qt__WA_TranslucentBackground, true)

	// window is layed out vertically
	w.WindowLayout = widgets.NewQVBoxLayout()
	w.WindowLayout.SetContentsMargins(0, 0, 0, 0)
	w.SetLayout(w.WindowLayout)

	// add a inner window widget (since window is completely transparent this is needed for the border)
	w.InnerWindow = widgets.NewQFrame(nil, 0)
	w.InnerWindow.SetObjectName("innerWindow")
	w.WindowLayout.AddWidget(w.InnerWindow, 0, 0)

	// inner window is layed vertically
	w.Layout = widgets.NewQVBoxLayout()
	w.Layout.SetContentsMargins(0, 0, 0, 0)
	w.InnerWindow.SetLayout(w.Layout)

	// LineEdit for the search query
	w.Input = widgets.NewQLineEdit(nil)
	w.Input.SetObjectName("input")
	w.Input.SetFixedHeight(EDITLINE_HEIGHT)
	w.Layout.AddWidget(w.Input, 0, 0)

	// list of logins
	w.List = widgets.NewQListView(nil)
	w.List.SetObjectName("list")
	// expands horizontally but sticks to size hint vertically
	w.List.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Fixed)
	// hidden by default (shown as soon as something appears inside)
	w.List.Hide()

	// items in the list of logins
	w.Item = widgets.NewQStyledItemDelegate(nil)
	// set the list items
	w.List.SetItemDelegate(w.Item)
	w.Layout.AddWidget(w.List, 0, 0)

	w.setupEventListeners()
	return w
}

func (w *SearchUI) setupEventListeners() {
	// exit on esc
	w.Input.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		if event.Key() == int(core.Qt__Key_Escape) {
			ui.Quit()
		} else {
			w.Input.KeyPressEventDefault(event)
		}
	})

	w.Item.ConnectSizeHint(func(option *widgets.QStyleOptionViewItem, index *core.QModelIndex) *core.QSize {
		return core.NewQSize2(w.Item.SizeHintDefault(option, index).Width(), RESULT_HEIGHT)
	})
	w.List.ConnectMinimumSizeHint(func() *core.QSize {
		return core.NewQSize2(0, 0)
	})
	// size of the list (don't let the list grow bigger than 5 items)
	w.List.ConnectSizeHint(func() *core.QSize {
		rowSize := w.List.SizeHintForRow(0)
		// -1 when hidden should be 0
		if rowSize < 0 {
			rowSize = 0
		}
		count := w.List.Model().RowCount(core.NewQModelIndex())
		if count > 5 {
			count = 5
		}
		return core.NewQSize2(w.List.SizeHintDefault().Width(), RESULT_HEIGHT*count)
	})

	// make window draggable
	var xOffset, yOffset int
	w.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		xOffset = event.X()
		yOffset = event.Y()
	})
	w.ConnectMouseMoveEvent(func(event *gui.QMouseEvent) {
		w.Move2(event.GlobalX()-xOffset, event.GlobalY()-yOffset)
	})

	w.Input.ConnectTextChanged(func(text string) {
		filter(text)
	})
}

// updates the size of the list (auto hides if list model is empty)
func (w *SearchUI) UpdateSize() {
	count := w.List.Model().RowCount(core.NewQModelIndex())
	// hide if no items in the list
	if count > 0 {
		w.List.Show()
		w.List.UpdateGeometry()
	} else {
		w.List.Hide()
	}

	// forces parents to resize in case list has become smaller
	parent := w.List.ParentWidget()
	for parent != nil {
		parent.AdjustSize()
		parent = parent.ParentWidget()
	}
}
