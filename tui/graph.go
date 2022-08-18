package tui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func NewPlot() (plot *widgets.Plot) {
	plot = widgets.NewPlot()
	plot.Title = "Plot"
	plot.Marker = widgets.MarkerDot
	plot.Data = make([][]float64, 2)
	// plot.SetRect(50, 0, 75, 10)
	plot.DotMarkerRune = '+'
	plot.AxesColor = ui.ColorWhite
	plot.LineColors[0] = ui.ColorYellow
	plot.DrawDirection = widgets.DrawLeft

	return
}
