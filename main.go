package main

import (
	"log"
	"time"

	"miniminiz/server"
	"miniminiz/tui"

	ui "github.com/gizak/termui/v3"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	valueBuffer := make(chan server.ActiveItemValue, 100)

	go server.RunServer("127.0.0.1", "10051", valueBuffer)

	// var valueStore map[string]*cache.HostValueCache
	valueStore := make(map[string]*server.HostValueCache)

	hosts := server.Monitoring.GetHosts()
	for name, items := range *hosts {
		valueStore[name] = server.NewHostValueCache(items)
	}

	go func() {
		for item := range valueBuffer {
			valueStore[item.Host].SaveValue(item)
		}
	}()

	hostlist := tui.NewHostList([]string{"localhost"})

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	// r := ui.NewRow(0.5, p1)
	table1 := tui.NewTable("system.cpu.load")
	table2 := tui.NewTable("system.sw.os")

	grid.Set(
		ui.NewRow(1.0,
			ui.NewCol(0.1,
				ui.NewRow(0.5, hostlist),
			),
			ui.NewCol(0.9,
				ui.NewRow(0.5, table1),
				ui.NewRow(0.5, table2),
			),
		),
	)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>", "e":
				return
			}
		case <-ticker:

			vals := valueStore["localhost"].GetValues("system.cpu.load")
			tui.UpdateTable(table1, vals)

			vals2 := valueStore["localhost"].GetValues("system.sw.os")
			tui.UpdateTable(table2, vals2)

			ui.Render(grid)
		}
	}

}
