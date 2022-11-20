package main

import (
	"log"
	"os"
	"strings"
	"time"

	"miniminiz/server"
	"miniminiz/tui"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	configFile := os.Args[1]

	t, _ := server.NewServer(configFile)

	go t.RunServer()

	valueStore := make(map[string]*server.HostValueCache)

	for name, items := range *t.Monitoring.GetHosts() {
		valueStore[name] = server.NewHostValueCache(items)
	}

	go func() {
		for item := range t.Cache {
			valueStore[item.Host].SaveValue(item)
		}
	}()

	hostlist := tui.NewList("Hosts", t.Monitoring.GetHostnames())

	a := t.Monitoring.GetKeys(hostlist.Rows[hostlist.SelectedRow])

	useri := tui.NewUserInterface()
	useri.HostList = hostlist
	useri.MetricList = tui.NewList("Metrics", a)
	useri.SetGrid()

	ui.Render(useri.Grid)
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C

	ui.Render(useri.Grid)
	activeList := useri.HostList
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "j", "<Down>":
				activeList.ScrollDown()
				if activeList == useri.HostList {
					useri.MetricList.Rows = t.Monitoring.GetKeys(hostlist.Rows[hostlist.SelectedRow])
				}
			case "k", "<Up>":
				activeList.ScrollUp()
				if activeList == useri.HostList {
					useri.MetricList.Rows = t.Monitoring.GetKeys(hostlist.Rows[hostlist.SelectedRow])
				}
			case "<Enter>":
				if activeList == useri.HostList {
					activeList = useri.MetricList
				}
			case "t":
				if activeList == useri.MetricList {
					currentHost := hostlist.Rows[hostlist.SelectedRow]
					currentMetric := useri.MetricList.Rows[useri.MetricList.SelectedRow]
					useri.BottomWidget = useri.WidgetCache.BottomTable
					tui.RenameTable(useri.BottomWidget.(*widgets.Table), currentHost, currentMetric)
				}
			case "T":
				if activeList == useri.MetricList {
					currentHost := hostlist.Rows[hostlist.SelectedRow]
					currentMetric := useri.MetricList.Rows[useri.MetricList.SelectedRow]
					useri.TopWidget = useri.WidgetCache.TopTable
					tui.RenameTable(useri.TopWidget.(*widgets.Table), currentHost, currentMetric)
				}
			case "<Escape>":
				activeList = useri.HostList
			case "q", "<C-c>", "e":
				return
			}
		case <-ticker:
			top := useri.TopWidget.(*widgets.Table)
			if top.Title != "Dummy:Dummy" {
				topSplit := strings.Split(top.Title, ":")
				topValues := valueStore[topSplit[0]].GetValues(topSplit[1])
				tui.UpdateTable(top, topValues)
			}
			bottom := useri.BottomWidget.(*widgets.Table)
			if bottom.Title != "Dummy:Dummy" {
				bottomSplit := strings.Split(bottom.Title, ":")
				bottomValues := valueStore[bottomSplit[0]].GetValues(bottomSplit[1])
				tui.UpdateTable(bottom, bottomValues)
			}
		}
		useri.SetGrid()
		ui.Render(useri.Grid)
	}

}
