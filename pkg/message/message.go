package message

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
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

// StringFreq is a struct with a Value and its Frequency
type StringFreq struct {
	Value string
	Freq  int
}

// SortedAnalysis contains the aggregate analysis with the fields sorted
type SortedAnalysis struct {
	ParticipantAnalyses []*SortedParticipantAnalysis
	Stickers            []StringFreq
	Words               []StringFreq
	Reactions           []StringFreq
	Mentions            []StringFreq
	MessageCount        int
}

// SortedParticipantAnalysis contains the sorted values of participants
type SortedParticipantAnalysis struct {
	Stickers     []StringFreq
	Words        []StringFreq
	Reactions    []StringFreq
	Mentions     []StringFreq
	MessageCount int
}

// Analysis contains the aggregate analysis and a map
// for participant to participant anaylsis
type Analysis struct {
	ParticipantAnalyses map[string]*ParticipantAnalysis
	Stickers            map[string]int
	Words               map[string]int
	Reactions           map[string]int
	Mentions            map[string]int
	MessageCount        int
}

// ParticipantAnalysis contains the participant analysis for
// sticker, word, reactions
type ParticipantAnalysis struct {
	Stickers     map[string]int
	Words        map[string]int
	Reactions    map[string]int
	Mentions     map[string]int
	MessageCount int
}

func newParticipantAnalysis() *ParticipantAnalysis {
	return &ParticipantAnalysis{
		Stickers:     make(map[string]int),
		Words:        make(map[string]int),
		Reactions:    make(map[string]int),
		Mentions:     make(map[string]int),
		MessageCount: 0,
	}
}

func newAnalysis() Analysis {
	return Analysis{
		ParticipantAnalyses: make(map[string]*ParticipantAnalysis),
		Stickers:            make(map[string]int),
		Words:               make(map[string]int),
		Reactions:           make(map[string]int),
		Mentions:            make(map[string]int),
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
		err := AnalyzeMessage(&a, m)
		if err != nil {
			fmt.Printf("analyzing message failed: %v", err)
		}
	}

	return a
}

// AnalyzeMessage processes a single message and changes the analysis
func AnalyzeMessage(a *Analysis, m Message) error {
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
		return nil
	}
	if strings.Contains(m.Content, "sent an attachment.") {
		return nil
	}
	if strings.Contains(m.Content, "sent a sticker.") {
		if _, ok := a.Stickers[m.Sticker.URI]; ok {
			a.Stickers[m.Sticker.URI]++
		} else {
			a.Stickers[m.Sticker.URI] = 1
		}

		if _, ok := a.ParticipantAnalyses[m.SenderName].Stickers[m.Sticker.URI]; ok {
			a.ParticipantAnalyses[m.SenderName].Stickers[m.Sticker.URI]++
		} else {
			a.ParticipantAnalyses[m.SenderName].Stickers[m.Sticker.URI] = 1
		}

		return nil
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9@]+")
	if err != nil {
		return errors.Wrap(err, "regex failed to compile")
	}
	words := reg.Split(strings.ToLower(m.Content), -1)
	for _, word := range words {
		if len(word) == 0 {
			continue
		}

		if word[0] == '@' && len(word) > 1 {
			if _, ok := a.Mentions[word]; ok {
				a.Mentions[word]++
			} else {
				a.Mentions[word] = 1
			}

			if _, ok := a.ParticipantAnalyses[m.SenderName].Mentions[word]; ok {
				a.ParticipantAnalyses[m.SenderName].Mentions[word]++
			} else {
				a.ParticipantAnalyses[m.SenderName].Mentions[word] = 1
			}
		}

		if _, ok := a.Words[word]; ok {
			a.Words[word]++
		} else {
			a.Words[word] = 1
		}

		if _, ok := a.ParticipantAnalyses[m.SenderName].Words[word]; ok {
			a.ParticipantAnalyses[m.SenderName].Words[word]++
		} else {
			a.ParticipantAnalyses[m.SenderName].Words[word] = 1
		}
	}

	return nil
}

// SortAnalysis returns the sorted analysis which each map sorted
func SortAnalysis(a Analysis) SortedAnalysis {
	s := SortedAnalysis{}

	return s
}
