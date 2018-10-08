package message

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Blob is the struct to represent the entire message.json file
type Blob struct {
	Participants []Participant `json:"participants"`
	Messages     []Message     `json:"messages"`
}

// Participant is the struct to represent each participant
type Participant struct {
	Name string `json:"name"`
}

// Message is the struct to represent each message
type Message struct {
	SenderName  string      `json:"sender_name"`
	TimestampMs int64       `json:"timestamp_ms"`
	Content     string      `json:"content"`
	Sticker     *Sticker    `json:"sticker"`
	Reactions   *[]Reaction `json:"reactions"`
	Type        string      `json:"type"`
}

// Sticker is the struct to represent the message sticker
type Sticker struct {
	URI string `json:"uri"`
}

// Reaction is the struct to represent the message reaction
type Reaction struct {
	Reaction string `json:"reaction"`
	Actor    string `json:"actor"`
}

// ParseMessages returns the data structure for the message.json
func ParseMessages(filepath string) (Blob, error) {
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Blob{}, errors.Wrap(err, "failed to read file")
	}

	var b Blob
	err = json.Unmarshal(dat, &b)
	if err != nil {
		return Blob{}, errors.Wrap(err, "json failed to parse")
	}

	return b, nil
}
