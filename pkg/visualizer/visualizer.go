package visualizer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kevinwubert/fb-messenger-analysis/pkg/message"
	"github.com/pkg/errors"

	"github.com/wcharczuk/go-chart"
)

type client struct {
	SortedAnalysis message.SortedAnalysis
}

// Client returns a client for the visualizer
type Client interface {
	DrawBarGraphHandler(w http.ResponseWriter, r *http.Request)
	GetNamesHandler(w http.ResponseWriter, r *http.Request)
	TopStickerHandler(w http.ResponseWriter, r *http.Request)
}

// New returns a new client for the visualizer
func New(sortedAnalysis message.SortedAnalysis) Client {
	return client{
		SortedAnalysis: sortedAnalysis,
	}
}

func (c client) DrawBarGraphHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if _, ok := query["name"]; !ok {
		fmt.Printf("no name query")
		WriteErrorResponse(w, errors.New("no name query"))
		return
	}
	if _, ok := query["type"]; !ok {
		fmt.Printf("no type query")
		WriteErrorResponse(w, errors.New("no type query"))
		return
	}
	if _, ok := query["count"]; !ok {
		fmt.Printf("no count query")
		WriteErrorResponse(w, errors.New("no count query"))
		return
	}

	name := query["name"][0]
	queryType := query["type"][0]
	countStr := query["count"][0]

	count, err := strconv.Atoi(countStr)
	if err != nil {
		fmt.Println("failed to parse count: %v", err)
		WriteErrorResponse(w, errors.Wrap(err, "failed to parse count"))
		return
	}

	if name != "everyone" {
		if _, ok := c.SortedAnalysis.SortedParticipantAnalyses[name]; !ok {
			fmt.Printf("invalid name")
			WriteErrorResponse(w, errors.New("invalid name"))
			return
		}
	}

	bc := chart.BarChart{
		Title:      GetGraphTitle(name, queryType, countStr),
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		XAxis: chart.Style{
			Show: true,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Bars: c.GetValuesFromQuery(name, queryType, count),
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err = bc.Render(chart.PNG, w)

	if err != nil {
		fmt.Printf("Error rendering pie chart: %v\n", err)
	}
}

// GetGraphTitle gets the graph title
func GetGraphTitle(name string, queryType string, count string) string {
	return "Top " + count + " " + queryType + " for " + name
}

// GetValuesFromQuery gets the values for the bar graph
func (c client) GetValuesFromQuery(name string, queryType string, count int) []chart.Value {
	values := []chart.Value{}

	if name == "everyone" {
		switch queryType {
		case "words":
			for _, v := range c.SortedAnalysis.Words {
				values = append(values, chart.Value{Value: float64(v.Freq), Label: v.Value})
			}
		case "stickers":
			for _, v := range c.SortedAnalysis.Stickers {
				values = append(values, chart.Value{Value: float64(v.Freq), Label: v.Value})
			}
		case "mentions":
			for _, v := range c.SortedAnalysis.Mentions {
				values = append(values, chart.Value{Value: float64(v.Freq), Label: v.Value})
			}
		case "reactions":
			for _, v := range c.SortedAnalysis.Reactions {
				values = append(values, chart.Value{Value: float64(v.Freq), Label: v.Value})
			}
		}
	} else {
		switch queryType {
		case "words":
			for _, v := range c.SortedAnalysis.SortedParticipantAnalyses[name].Words {
				values = append(values, chart.Value{Value: float64(v.Freq), Label: v.Value})
			}
		case "stickers":
			for _, v := range c.SortedAnalysis.SortedParticipantAnalyses[name].Stickers {
				values = append(values, chart.Value{Value: float64(v.Freq), Label: v.Value})
			}
		case "mentions":
			for _, v := range c.SortedAnalysis.SortedParticipantAnalyses[name].Mentions {
				values = append(values, chart.Value{Value: float64(v.Freq), Label: v.Value})
			}
		case "reactions":
			for _, v := range c.SortedAnalysis.SortedParticipantAnalyses[name].Reactions {
				values = append(values, chart.Value{Value: float64(v.Freq), Label: v.Value})
			}
		}
	}
	if len(values) > count {
		return values[:count]
	}
	return values
}

func (c client) TopStickerHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if _, ok := query["name"]; !ok {
		fmt.Printf("no name query")
		WriteErrorResponse(w, errors.New("no name query"))
		return
	}

	name := query["name"][0]

	if name != "everyone" {
		if _, ok := c.SortedAnalysis.SortedParticipantAnalyses[name]; !ok {
			fmt.Printf("invalid name")
			WriteErrorResponse(w, errors.New("invalid name"))
			return
		}
	}

	var stickerID string
	if name == "everyone" {
		stickerID = c.SortedAnalysis.Stickers[0].Value
	} else {
		stickerID = c.SortedAnalysis.SortedParticipantAnalyses[name].Stickers[0].Value
	}

	url := "https://messenger.com/stickers/asset/?sticker_id=" + stickerID
	fmt.Println(url)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (c client) GetNamesHandler(w http.ResponseWriter, r *http.Request) {
	names := []string{"everyone"}

	for k := range c.SortedAnalysis.SortedParticipantAnalyses {
		names = append(names, k)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	WriteJSONResponse(w, names)
}

// WriteErrorResponse writes an error back from an invalid request
func WriteErrorResponse(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// WriteJSONResponse writes some value and encodes into the response
func WriteJSONResponse(w http.ResponseWriter, val interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(val)
}
