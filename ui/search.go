package ui

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const (
	appName           = "1Password Lookup"
	searchWidth       = 500
	searchInputHeight = 50
	resultHeight      = 50
	maxResults        = 5
)

const searchStyles = `
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

// SearchUI is a "class" that is the search popup
// it stores the Qt objects and provides the slots to manipulate it
type SearchUI struct {
	widgets.QWidget

	_ func()       `constructor:"init"`
	_ func()       `slot:"Start"` // start/finish because show/hide would collide with QWidget's show/hide
	_ func()       `slot:"Finish"`
	_ func()       `slot:"UpdateSize"`
	_ func()       `slot:"ListDataWillChange`
	_ func()       `slot:"ListDataDidChange`
	_ func()       `slot:"Disable`
	_ func()       `slot:"EnableAndFocus`
	_ func(string) `slot:"Copy`

	windowLayout *widgets.QVBoxLayout
	innerWindow  *widgets.QFrame
	layout       *widgets.QVBoxLayout
	input        *widgets.QLineEdit
	item         *widgets.QStyledItemDelegate
	list         *widgets.QListView
	listModel    *core.QAbstractListModel

	listData    func(int, int) string
	listCount   func() int
	textChanged func(string)
	copyItem    func(int)
	openItem    func(int)
}

func (w *SearchUI) init() {
	// create the window
	w.SetWindowTitle(appName)
	w.SetStyleSheet(searchStyles)

	// initially in waiting state
	cursor := gui.NewQCursor2(core.Qt__WaitCursor)
	w.SetCursor(cursor)

	// tell window to quit when it closes (Qt::Tool turns this off for some reason)
	w.SetAttribute(core.Qt__WA_QuitOnClose, true)
	w.SetAttribute(core.Qt__WA_TranslucentBackground, true)

	// window is layed out vertically
	w.windowLayout = widgets.NewQVBoxLayout()
	w.windowLayout.SetContentsMargins(0, 0, 0, 0)
	w.SetLayout(w.windowLayout)

	// add a inner window widget (since window is completely transparent this is needed for the border)
	w.innerWindow = widgets.NewQFrame(nil, 0)
	w.innerWindow.SetObjectName("innerWindow")
	w.windowLayout.AddWidget(w.innerWindow, 0, 0)

	// inner window is layed vertically
	w.layout = widgets.NewQVBoxLayout()
	w.layout.SetContentsMargins(0, 0, 0, 0)
	w.innerWindow.SetLayout(w.layout)

	// LineEdit for the search query
	w.input = widgets.NewQLineEdit(nil)
	w.input.SetObjectName("input")
	w.input.SetDisabled(true)
	w.input.SetFixedHeight(searchInputHeight)
	w.layout.AddWidget(w.input, 0, 0)

	// list of logins
	w.list = widgets.NewQListView(nil)
	w.list.SetObjectName("list")
	// expands horizontally but sticks to size hint vertically
	w.list.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Fixed)
	// hidden by default (shown as soon as something appears inside)
	w.list.Hide()

	// items in the list of logins
	w.item = widgets.NewQStyledItemDelegate(nil)
	// set the list items
	w.list.SetItemDelegate(w.item)
	w.layout.AddWidget(w.list, 0, 0)

	// model for the list to provide data
	w.listModel = core.NewQAbstractListModel(nil)
	w.list.SetModel(w.listModel)

	// set size hint
	w.input.ConnectSizeHint(func() *core.QSize {
		return core.NewQSize2(searchWidth, w.input.SizeHintDefault().Height())
	})
	w.item.ConnectSizeHint(func(option *widgets.QStyleOptionViewItem, index *core.QModelIndex) *core.QSize {
		return core.NewQSize2(w.item.SizeHintDefault(option, index).Width(), resultHeight)
	})
	w.list.ConnectMinimumSizeHint(func() *core.QSize {
		return core.NewQSize2(0, 0)
	})
	w.list.ConnectSizeHint(w.listSize)

	// make window draggable
	var xOffset, yOffset int
	w.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		xOffset = event.X()
		yOffset = event.Y()
	})
	w.ConnectMouseMoveEvent(func(event *gui.QMouseEvent) {
		w.Move2(event.GlobalX()-xOffset, event.GlobalY()-yOffset)
	})

	w.input.ConnectKeyPressEvent(w.keyListener)
	w.input.ConnectReturnPressed(w.returnListener)
	w.input.ConnectTextChanged(func(text string) {
		if w.textChanged != nil {
			w.textChanged(text)
		}
	})

	// focus on input when enabled
	w.input.ConnectChangeEvent(func(event *core.QEvent) {
		if event.Type() == core.QEvent__EnabledChange {
			w.input.SetFocus(core.Qt__NoFocusReason)
		} else {
			w.input.ChangeEventDefault(event)
		}
	})

	w.listModel.ConnectFlags(func(index *core.QModelIndex) core.Qt__ItemFlag {
		return core.Qt__ItemIsSelectable | core.Qt__ItemIsEnabled
	})
	// setup list data
	w.listModel.ConnectData(func(index *core.QModelIndex, role int) *core.QVariant {
		if w.listData != nil {
			if strData := w.listData(index.Row(), role); strData != "" {
				return core.NewQVariant17(strData)
			}
		}
		return core.NewQVariant()
	})
	w.listModel.ConnectRowCount(func(parent *core.QModelIndex) int {
		if w.listCount != nil {
			return w.listCount()
		}
		return 0
	})

	// setup slots
	w.ConnectStart(func() {
		w.Show()
	})
	w.ConnectFinish(func() {
		w.Hide()
	})
	w.ConnectUpdateSize(w.updateSize)
	w.ConnectListDataWillChange(func() {
		w.listModel.LayoutAboutToBeChanged(nil, core.QAbstractItemModel__NoLayoutChangeHint)
	})
	w.ConnectListDataDidChange(func() {
		w.listModel.LayoutChanged(nil, core.QAbstractItemModel__NoLayoutChangeHint)
	})
	w.ConnectDisable(w.disable)
	w.ConnectEnableAndFocus(w.enableAndFocus)
	w.ConnectCopy(func(value string) {
		App.Clipboard().SetText(value, gui.QClipboard__Clipboard)
	})
}

func (w *SearchUI) SetListDataProviders(data func(int, int) string, count func() int, copy func(int), open func(int)) {
	w.listData = data
	w.listCount = count
	w.copyItem = copy
	w.openItem = open
}

func (w *SearchUI) SetTextChanged(f func(string)) {
	w.textChanged = f
}

// keyListener will listen for special keys to handle specially, otherwise will call default event handler
// esc will quit, up/down will move list selection
func (w *SearchUI) keyListener(event *gui.QKeyEvent) {
	if event.Key() == int(core.Qt__Key_Escape) {
		App.Quit()
	} else if event.Key() == int(core.Qt__Key_Down) || event.Key() == int(core.Qt__Key_Up) {
		row := w.list.CurrentIndex().Row()
		if event.Key() == int(core.Qt__Key_Down) {
			row += 1
		} else if event.Key() == int(core.Qt__Key_Up) {
			row -= 1
		}
		rowIndex := w.listModel.Index(row, 0, w.list.RootIndex())
		// valid means the index isn't out of range
		if rowIndex.IsValid() {
			w.list.SetCurrentIndex(rowIndex)
		}
	} else if event.Key() == int(core.Qt__Key_C) && (event.Modifiers()&core.Qt__ControlModifier != 0) {
		indexes := w.list.SelectedIndexes()
		if len(indexes) > 0 {
			selected := indexes[0].Row()
			w.copyItem(selected)
		}
	} else {
		w.input.KeyPressEventDefault(event)
	}
}

// returnListener will respond to enter key, opens details for an item
func (w *SearchUI) returnListener() {
	indexes := w.list.SelectedIndexes()
	if len(indexes) > 0 {
		selected := indexes[0].Row()
		w.openItem(selected)
	}
}

// listSize will calcualte the size hint of the list (height only) based on number of items
// this will cap the size of the list to a maximum number of items
func (w *SearchUI) listSize() *core.QSize {
	rowSize := w.list.SizeHintForRow(0)
	// -1 when hidden should be 0
	if rowSize < 0 {
		rowSize = 0
	}
	count := w.list.Model().RowCount(core.NewQModelIndex())
	if count > maxResults {
		count = maxResults
	}
	return core.NewQSize2(w.list.SizeHintDefault().Width(), resultHeight*count)
}

// updateSize updates the size of the list (auto hides if list model is empty)
func (w *SearchUI) updateSize() {
	count := w.list.Model().RowCount(core.NewQModelIndex())
	// hide if no items in the list
	if count > 0 {
		w.list.Show()
		w.list.UpdateGeometry()
	} else {
		w.list.Hide()
	}

	// forces parents to resize in case list has become smaller
	parent := w.list.ParentWidget()
	for parent != nil {
		parent.AdjustSize()
		parent = parent.ParentWidget()
	}
}

// disable will disable the input
func (w *SearchUI) disable() {
	w.input.SetEnabled(false)
	cursor := gui.NewQCursor2(core.Qt__WaitCursor)
	w.SetCursor(cursor)
}

// enableAndFocus will enable the input which will trigger a
// listener that wll focus the input when it is enabled
func (w *SearchUI) enableAndFocus() {
	w.input.SetEnabled(true)
	cursor := gui.NewQCursor2(core.Qt__ArrowCursor)
	w.SetCursor(cursor)
}
