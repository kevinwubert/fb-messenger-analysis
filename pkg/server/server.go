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
	fmt.Println("analyzing messages...")
	analysis := message.AnalyzeMessages(messageBlob)
	sortedAnalysis := message.SortAnalysis(analysis)
	fmt.Println("finished analyzing messages...")
	fmt.Println("starting facebook messenger analysis server...")

	visualizerClient := visualizer.New(sortedAnalysis)
	http.HandleFunc("/graph", visualizerClient.DrawBarGraphHandler)
	http.HandleFunc("/topSticker", visualizerClient.TopStickerHandler)
	http.HandleFunc("/getNames", visualizerClient.GetNamesHandler)

	err = http.ListenAndServe(":80", nil)
	return err
}
