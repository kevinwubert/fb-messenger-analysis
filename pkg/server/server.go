package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kevinwubert/fb-messenger-analysis/pkg/message"
	"github.com/kevinwubert/fb-messenger-analysis/pkg/visualizer"
	"github.com/pkg/errors"
)

func Main() error {
	args := os.Args
	if len(args) != 2 {
		return errors.New("invalid number of arguments.\ncommand format: fb-messenger-analysis <message.json filepath>")
	}

	messageFilepath := os.Args[1]
	messageBlob, err := message.ParseMessages(messageFilepath)
	if err != nil {
		return errors.Wrap(err, "failed to parse messages")
	}

	analysis := message.AnalyzeMessages(messageBlob)
	sortedAnalysis := message.SortAnalysis(analysis)
	fmt.Println(sortedAnalysis.SortedParticipantAnalyses["Kevin Wu"])

	fmt.Println("Starting Facebook Messenger Analysis server...")
	http.HandleFunc("/graph/", visualizer.DrawPieChartHandler)

	err = http.ListenAndServe(":8080", nil)
	return err
}
