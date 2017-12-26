package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const (
	appName        = "1Password Lookup"
	windowWidth    = 600
	editlineHeight = 50
	resultHeight   = 50
)

// SearchUI is a "class" that is the search popup
// it stores the Qt objects and provides the slots to manipulate it
type SearchUI struct {
	widgets.QWidget

	_ func()                          `slots:"Show"`
	_ func()                          `slots:"Hide"`
	_ func()                          `slots:"Disable"`
	_ func() *core.QAbstractListModel `slots:"GetNewListModel"`
	_ func(*core.QAbstractListModel)  `slots:"SetListModel"`
	_ func()                          `slots:"ListDataWillChange`
	_ func()                          `slots:"ListDataDidChange`
	_ func()                          `slots:"EnableAndFocus`

	WindowLayout *widgets.QVBoxLayout
	InnerWindow  *widgets.QFrame
	Layout       *widgets.QVBoxLayout
	Input        *widgets.QLineEdit
	Item         *widgets.QStyledItemDelegate
	List         *widgets.QListView
	ListModel    *core.QAbstractListModel
}

func setupSearch() *SearchUI {
	// create the window
	w := NewSearchUI(nil, core.Qt__Tool|
		core.Qt__FramelessWindowHint|
		core.Qt__WindowCloseButtonHint|
		core.Qt__WindowStaysOnTopHint)
	w.SetWindowTitle(appName)

	// initially in waiting state
	cursor := gui.NewQCursor2(core.Qt__WaitCursor)
	w.SetCursor(cursor)

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
	w.Input.SetDisabled(true)
	w.Input.SetFixedHeight(editlineHeight)
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

	// model for the list to provide data
	w.ListModel = core.NewQAbstractListModel(nil)
	w.ListModel.ConnectFlags(func(index *core.QModelIndex) core.Qt__ItemFlag {
		return core.Qt__ItemIsSelectable | core.Qt__ItemIsEnabled
	})
	w.List.SetModel(w.ListModel)

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
		return core.NewQSize2(w.Item.SizeHintDefault(option, index).Width(), resultHeight)
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
		return core.NewQSize2(w.List.SizeHintDefault().Width(), resultHeight*count)
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

	// focus on input when enabled
	w.Input.ConnectChangeEvent(func(event *core.QEvent) {
		if event.Type() == core.QEvent__EnabledChange {
			w.Input.SetFocus(core.Qt__NoFocusReason)
		} else {
			w.Input.ChangeEventDefault(event)
		}
	})
}

// UpdateSize updates the size of the list (auto hides if list model is empty)
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

// Show shows the search popup
func (w *SearchUI) Show() {
	w.QWidget.Show()
}

// Hide hides the search popup
func (w *SearchUI) Hide() {
	w.QWidget.Hide()
}

// SetupListModel will set a new model to list and setup two listeners provides for
// data and row count
func (w *SearchUI) SetupListModel(row func(int, int) string, count func() int) {
	w.ListDataWillChange()

	w.ListModel.ConnectData(func(index *core.QModelIndex, role int) *core.QVariant {
		strData := row(index.Row(), role)
		if strData != "" {
			return core.NewQVariant17(strData)
		}
		return core.NewQVariant()
	})
	w.ListModel.ConnectRowCount(func(parent *core.QModelIndex) int {
		return count()
	})

	w.ListDataDidChange()
}

// ListDataWillChange notifies the model data is going to change
func (w *SearchUI) ListDataWillChange() {
	w.ListModel.LayoutAboutToBeChanged(nil, core.QAbstractItemModel__NoLayoutChangeHint)
}

// ListDataDidChange notifies the model data has changed
func (w *SearchUI) ListDataDidChange() {
	w.ListModel.LayoutChanged(nil, core.QAbstractItemModel__NoLayoutChangeHint)
}

// EnableAndFocus will enable the input which will trigger a
// listener that wll focus the input when it is enabled
func (w *SearchUI) EnableAndFocus() {
	w.Input.SetDisabled(false)
	cursor := gui.NewQCursor2(core.Qt__ArrowCursor)
	w.SetCursor(cursor)
}
