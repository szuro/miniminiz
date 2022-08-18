package tui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func NewHostList(hosts []string) (hostlist *widgets.List) {
	hostlist = widgets.NewList()
	hostlist.Title = "Hosts"
	hostlist.Rows = hosts
	hostlist.TextStyle = ui.NewStyle(ui.ColorYellow)
	hostlist.WrapText = false
	hostlist.SetRect(0, 0, 25, 8)
	return
}
