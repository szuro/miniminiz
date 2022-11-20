package tui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type UserInterface struct {
	HostList     *widgets.List
	MetricList   *widgets.List
	TopWidget    interface{}
	BottomWidget interface{}
	Grid         *ui.Grid
	WidgetCache  widgetCache
}

type widgetCache struct {
	TopTable    *widgets.Table
	BottomTable *widgets.Table
}

func NewUserInterface() (userInt UserInterface) {
	userInt = UserInterface{}
	userInt.Grid = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	userInt.Grid.SetRect(0, 0, termWidth, termHeight)

	userInt.WidgetCache.TopTable = NewTable("Dummy", "Dummy")
	userInt.WidgetCache.BottomTable = NewTable("Dummy", "Dummy")

	userInt.TopWidget = userInt.WidgetCache.TopTable
	userInt.BottomWidget = userInt.WidgetCache.BottomTable

	return
}

func (useri *UserInterface) SetGrid() {
	useri.Grid.Set(
		ui.NewRow(0.5,
			ui.NewCol(0.2, useri.HostList),
			ui.NewCol(0.8, useri.TopWidget),
		),
		ui.NewRow(0.5,
			ui.NewCol(0.2, useri.MetricList),
			ui.NewCol(0.8, useri.BottomWidget),
		),
	)
}
