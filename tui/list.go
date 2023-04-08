package tui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func NewList(title string, entries []string) (hostlist *widgets.List) {
	hostlist = widgets.NewList()
	hostlist.Title = title
	hostlist.Rows = entries
	hostlist.TextStyle = ui.NewStyle(ui.ColorWhite)
	hostlist.SelectedRowStyle = ui.NewStyle(ui.ColorYellow)
	hostlist.WrapText = false
	hostlist.SetRect(0, 0, 25, 8)
	return
}
