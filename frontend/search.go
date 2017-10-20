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
	window       *widgets.QWidget
	windowLayout *widgets.QVBoxLayout
	innerWindow  *widgets.QFrame
	layout       *widgets.QVBoxLayout
	input        *widgets.QLineEdit
	item         *widgets.QStyledItemDelegate
	list         *widgets.QListView
}

func NewSearch() Search {
	fmt.Println("hi")

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

	// list of logins
	list := widgets.NewQListView(nil)
	list.SetObjectName("list")
	// expands horizontally but sticks to size hint vertically
	list.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Fixed)
	// hidden by default (shown as soon as something appears inside)
	list.Hide()

	// items in the list of logins
	item := widgets.NewQStyledItemDelegate(nil)
	// set the list items
	list.SetItemDelegate(item)
	layout.AddWidget(list, 0, 0)

	// exit on esc
	input.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		if event.Key() == int(core.Qt__Key_Escape) {
			CloseGui()
		} else {
			input.KeyPressEventDefault(event)
		}
	})

	// size of the list items
	item.ConnectSizeHint(func(option *widgets.QStyleOptionViewItem, index *core.QModelIndex) *core.QSize {
		// only care about setting the correct height
		return core.NewQSize2(item.SizeHintDefault(option, index).Width(), RESULT_HEIGHT)
	})
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
		count := list.Model().RowCount(core.NewQModelIndex())
		if count > 5 {
			count = 5
		}
		// only care about setting the correct height
		return core.NewQSize2(list.SizeHintDefault().Width(), RESULT_HEIGHT*count)
	})

	// make window draggable
	var xOffset, yOffset int
	window.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		xOffset = event.X()
		yOffset = event.Y()
	})
	window.ConnectMouseMoveEvent(func(event *gui.QMouseEvent) {
		window.Move2(event.GlobalX()-xOffset, event.GlobalY()-yOffset)
	})

	return Search{
		window,
		windowLayout,
		innerWindow,
		layout,
		input,
		item,
		list,
	}
}

func (s *Search) Show() {
	s.window.Show()
}

func (s *Search) Hide() {
	s.window.Hide()
}

func (s *Search) HideList() {
	s.list.Hide()
	s.updateParentSize()
}

func (s *Search) OnTextChanged(listener func(string)) {
	s.input.ConnectTextChanged(func(text string) {
		listener(text)
	})
}

func (s *Search) ContextMenuData(listener func(int) map[string]string) {
	// send a listener the
	s.list.ConnectContextMenuEvent(func(event *gui.QContextMenuEvent) {
		// get data
		selected := s.list.SelectedIndexes()[0].Row()
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
		menu.Exec3(actions, event.GlobalPos(), actions[0], s.list)
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

func (s *Search) SetListModel(model core.QAbstractItemModel_ITF) {
	s.list.SetModel(model)
}

func (s *Search) UpdateSize() {
	count := s.list.Model().RowCount(core.NewQModelIndex())
	// hide if empty
	if count > 0 {
		s.list.Show()
		s.list.UpdateGeometry()
	} else {
		s.list.Hide()
	}

	s.updateParentSize()
}

func (s *Search) updateParentSize() {
	parent := s.list.ParentWidget()
	for parent != nil {
		parent.AdjustSize()
		parent = parent.ParentWidget()
	}
}
