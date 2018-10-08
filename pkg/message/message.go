package message

import (
	"encoding/json"
	"io/ioutil"
	"strings"

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
// Ignoring sharing photo and attachments
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

// Analysis contains the aggregate analysis and a map
// for participant to participant anaylsis
type Analysis struct {
	ParticipantAnalyses map[string]*ParticipantAnalysis
	Stickers            map[string]int
	Words               map[string]int
	Reactions           map[string]int
	MessageCount        int
}

// ParticipantAnalysis contains the participant analysis for
// sticker, word, reactions
type ParticipantAnalysis struct {
	Stickers     map[string]int
	Words        map[string]int
	Reactions    map[string]int
	MessageCount int
	MentionCount int
}

func newParticipantAnalysis() *ParticipantAnalysis {
	return &ParticipantAnalysis{
		Stickers:     make(map[string]int),
		Words:        make(map[string]int),
		Reactions:    make(map[string]int),
		MessageCount: 0,
		MentionCount: 0,
	}
}

func newAnalysis() Analysis {
	return Analysis{
		ParticipantAnalyses: make(map[string]*ParticipantAnalysis),
		Stickers:            make(map[string]int),
		Words:               make(map[string]int),
		Reactions:           make(map[string]int),
		MessageCount:        0,
	}
}

// AnalyzeMessages analyzes the message blob and returns the results
func AnalyzeMessages(b Blob) Analysis {
	a := newAnalysis()
	for _, p := range b.Participants {
		a.ParticipantAnalyses[p.Name] = newParticipantAnalysis()
	}

	for _, m := range b.Messages {
		AnalyzeMessage(&a, m)
	}

	return a
}

// AnalyzeMessage processes a single message and changes the analysis
func AnalyzeMessage(a *Analysis, m Message) {
	if m.Reactions != nil {
		for _, r := range *m.Reactions {
			if _, ok := a.Reactions[r.Reaction]; ok {
				a.Reactions[r.Reaction]++
			} else {
				a.Reactions[r.Reaction] = 1
			}

			if _, ok := a.ParticipantAnalyses[r.Actor].Reactions[r.Reaction]; ok {
				a.ParticipantAnalyses[r.Actor].Reactions[r.Reaction]++
			} else {
				a.ParticipantAnalyses[r.Actor].Reactions[r.Reaction] = 1
			}
		}
	}

	a.MessageCount++
	a.ParticipantAnalyses[m.SenderName].MessageCount++

	if strings.Contains(m.Content, "sent a photo.") {
		return
	}
	if strings.Contains(m.Content, "sent an attachment.") {
		return
	}
	if strings.Contains(m.Content, "sent a sticker.") {
		a.Stickers()
		return
	}
}
