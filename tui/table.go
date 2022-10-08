package tui

import (
	"fmt"
	"miniminiz/server"

	"github.com/gizak/termui/v3/widgets"
)

func NewTable(host, metric string) (table *widgets.Table) {
	table = widgets.NewTable()
	table.Title = host + ": " + metric
	return
}

func UpdateTable(t *widgets.Table, metrics []server.ActiveItemValue) {
	for _, val := range metrics {
		time := fmt.Sprintf("%d", val.Clock)
		value := fmt.Sprintf("%f", val.Value)
		t.Rows = append(t.Rows, []string{time, value})
	}
}
