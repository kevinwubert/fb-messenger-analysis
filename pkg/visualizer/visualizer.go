package visualizer

import (
	"fmt"
	"net/http"

	"github.com/wcharczuk/go-chart"
)

func DrawPieChartHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Drawing Pie Chart")

	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: []chart.Value{
			{Value: 3, Label: "Blue"},
			{Value: 6, Label: "Red"},
		},
	}

	w.Header().Set("Content-Type", "image/png")
	err := pie.Render(chart.PNG, w)

	if err != nil {
		fmt.Printf("Error rendering pie chart: %v\n", err)
	}
}
