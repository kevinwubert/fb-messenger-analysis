package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kevinwubert/fb-messenger-analysis/pkg/visualizer"
)

func Main() error {
	args := os.Args
	if len(args) != 2 {
		return errors.New("invalid number of arguments.\ncommand format: fb-messenger-analysis <message.json filepath>")
	}

	messageFilepath := os.Args[1]
	fmt.Println(messageFilepath)

	fmt.Println("Starting Facebook Messenger Analysis server...")
	http.HandleFunc("/graph/", visualizer.DrawPieChartHandler)

	err := http.ListenAndServe(":8080", nil)
	return err
}
