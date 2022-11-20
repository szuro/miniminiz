package tui

import (
	"fmt"
	"miniminiz/server"
	"time"

	"github.com/gizak/termui/v3/widgets"
)

func NewTable(host, metric string) (table *widgets.Table) {
	table = widgets.NewTable()
	RenameTable(table, host, metric)
	table.Rows = append(table.Rows, []string{"Time", "Value"})
	return
}

func RenameTable(table *widgets.Table, host, metric string) {
	table.Title = fmt.Sprintf("%s:%s", host, metric)
}

func UpdateTable(t *widgets.Table, metrics []server.ActiveItemValue) {
	t.Rows = t.Rows[:1]
	//t_ := time.Now()
	// zone, offset := t_.Zone()
	for _, val := range metrics {
		tm := time.Unix(int64(val.Clock), int64(val.Ns))
		time := fmt.Sprintf("%s", tm.Format("2006-01-02 15:04:05"))
		value := fmt.Sprintf("%s", val.Value)
		t.Rows = append(t.Rows, []string{time, value})
	}
}
