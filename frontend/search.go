package frontend

import (
	"fmt"
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

type Search struct {
	widgets.QWidget

	_ func() `constructor:"init"`

	windowLayout *widgets.QVBoxLayout
	innerWindow  *widgets.QFrame
	layout       *widgets.QVBoxLayout
	input        *widgets.QLineEdit
	item         *widgets.QStyledItemDelegate
	list         *widgets.QListView
}

func (w *Search) init() {
	fmt.Println("hi")

	w.SetWindowTitle(APP_NAME)
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
	w.input.SetFixedHeight(EDITLINE_HEIGHT)
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

	// exit on esc
	w.input.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		if event.Key() == int(core.Qt__Key_Escape) {
			CloseGui()
		} else {
			w.input.KeyPressEventDefault(event)
		}
	})

	// size of the list items
	w.item.ConnectSizeHint(func(option *widgets.QStyleOptionViewItem, index *core.QModelIndex) *core.QSize {
		// only care about setting the correct height
		return core.NewQSize2(w.item.SizeHintDefault(option, index).Width(), RESULT_HEIGHT)
	})
	// let the list be empty (shrink to 0px)
	w.list.ConnectMinimumSizeHint(func() *core.QSize {
		return core.NewQSize2(0, 0)
	})
	// size of the list (don't let the list grow bigger than 5 items)
	w.list.ConnectSizeHint(func() *core.QSize {
		rowSize := w.list.SizeHintForRow(0)
		// -1 when hidden should be 0
		if rowSize < 0 {
			rowSize = 0
		}
		count := w.list.Model().RowCount(core.NewQModelIndex())
		if count > 5 {
			count = 5
		}
		// only care about setting the correct height
		return core.NewQSize2(w.list.SizeHintDefault().Width(), RESULT_HEIGHT*count)
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
}

func (w *Search) HideList() {
	w.list.Hide()
	w.updateParentSize()
}

func (w *Search) OnTextChanged(listener func(string)) {
	w.input.ConnectTextChanged(func(text string) {
		listener(text)
	})
}

func (w *Search) ContextMenuData(listener func(int) map[string]string) {
	// send a listener the
	w.list.ConnectContextMenuEvent(func(event *gui.QContextMenuEvent) {
		// get data
		selected := w.list.SelectedIndexes()[0].Row()
		data := listener(selected)

		if data == nil {
			// no data probably logged out
			return
		}

		// each data item is a copy action
		actions := make([]*widgets.QAction, 0)
		for key, value := range data {
			action := createCopyAction(key, value)
			actions = append(actions, action)
		}

		menu := widgets.NewQMenu(nil)
		menu.Exec3(actions, event.GlobalPos(), actions[0], w.list)
	})
}

func createCopyAction(key, value string) *widgets.QAction {
	label := fmt.Sprintf("Copy %v", key)
	action := widgets.NewQAction2(label, nil)
	action.ConnectTriggered(func(checked bool) {
		app.Clipboard().SetText(value, gui.QClipboard__Clipboard)
	})
	return action
}

func (w *Search) SetListModel(model core.QAbstractItemModel_ITF) {
	w.list.SetModel(model)
}

func (w *Search) UpdateSize() {
	count := w.list.Model().RowCount(core.NewQModelIndex())
	// hide if empty
	if count > 0 {
		w.list.Show()
		w.list.UpdateGeometry()
	} else {
		w.list.Hide()
	}

	w.updateParentSize()
}

func (w *Search) updateParentSize() {
	parent := w.list.ParentWidget()
	for parent != nil {
		parent.AdjustSize()
		parent = parent.ParentWidget()
	}
}
