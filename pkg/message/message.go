package message

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

var ignoreWords = []string{"a", "about", "above", "after", "again", "against",
	"all", "am", "an", "and", "any", "are", "aren't", "as", "at",
	"be", "because", "been", "before", "being", "below", "between",
	"both", "but", "by", "can't", "cannot", "could", "couldn't",
	"did", "didn't", "do", "does", "doesn't", "doing", "don't",
	"down", "during", "each", "few", "for", "from", "further", "had",
	"hadn't", "has", "hasn't", "have", "haven't", "having", "he",
	"he'd", "he'll", "he's", "her", "here", "here's", "hers", "herself",
	"him", "himself", "his", "how", "how's", "i", "i'd", "i'll", "i'm",
	"i've", "if", "in", "into", "is", "isn't", "it", "it's", "its",
	"itself", "let's", "me", "more", "most", "mustn't", "my", "myself",
	"no", "nor", "not", "of", "off", "on", "once", "only", "or", "other",
	"ought", "our", "ours  ourselves", "out", "over", "own", "same",
	"shan't", "she", "she'd", "she'll", "she's", "should", "shouldn't",
	"so", "some", "such", "than", "that", "that's", "the", "their",
	"theirs", "them", "themselves", "then", "there", "there's", "these",
	"they", "they'd", "they'll", "they're", "they've", "this", "those",
	"through", "to", "too", "under", "until", "up", "very", "was", "wasn't",
	"we", "we'd", "we'll", "we're", "we've", "were", "weren't", "what",
	"what's", "when", "when's", "where", "where's", "which", "while",
	"who", "who's", "whom", "why", "why's", "with", "won't", "would",
	"wouldn't", "you", "you'd", "you'll", "you're", "you've", "your",
	"yours", "yourself", "yourselves", "http", "https", "www", "com",
	":", "http:", "https:", "im"}

func stringsToMap(strings []string) map[string]bool {
	m := make(map[string]bool)

	for _, s := range strings {
		m[s] = true
	}

	return m
}

var ignoreWordMap = stringsToMap(ignoreWords)

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

func nameToFirstName(s string) string {
	firstName := strings.Split(s, " ")[0]
	return firstName
}

// AnalyzeMessages analyzes the message blob and returns the results
func AnalyzeMessages(b Blob) Analysis {
	a := newAnalysis()
	for _, p := range b.Participants {
		a.ParticipantAnalyses[nameToFirstName(p.Name)] = newParticipantAnalysis()
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

			if _, ok := a.ParticipantAnalyses[nameToFirstName(r.Actor)].Reactions[r.Reaction]; ok {
				a.ParticipantAnalyses[nameToFirstName(r.Actor)].Reactions[r.Reaction]++
			} else {
				a.ParticipantAnalyses[nameToFirstName(r.Actor)].Reactions[r.Reaction] = 1
			}
		}
	}

	a.MessageCount++
	a.ParticipantAnalyses[nameToFirstName(m.SenderName)].MessageCount++

	if strings.Contains(m.Content, "sent a photo.") {
		return nil
	}
	if strings.Contains(m.Content, "sent an attachment.") {
		return nil
	}
	if m.Sticker != nil {
		stickerID := strings.Split(m.Sticker.URI, "_n_")[1]
		stickerID = strings.Split(stickerID, ".")[0]

		if stickerID == "369239263222822" {
			return nil
		}

		if _, ok := a.Stickers[stickerID]; ok {
			a.Stickers[stickerID]++
		} else {
			a.Stickers[stickerID] = 1
		}

		if _, ok := a.ParticipantAnalyses[nameToFirstName(m.SenderName)].Stickers[stickerID]; ok {
			a.ParticipantAnalyses[nameToFirstName(m.SenderName)].Stickers[stickerID]++
		} else {
			a.ParticipantAnalyses[nameToFirstName(m.SenderName)].Stickers[stickerID] = 1
		}

		return nil
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9@':]+")
	if err != nil {
		return errors.Wrap(err, "regex failed to compile")
	}
	words := reg.Split(strings.ToLower(m.Content), -1)
	for _, word := range words {
		if len(word) <= 1 {
			continue
		}

		if _, ok := ignoreWordMap[word]; ok {
			continue
		}

		if word[0] == '@' && len(word) > 1 {
			if _, ok := a.Mentions[word]; ok {
				a.Mentions[word]++
			} else {
				a.Mentions[word] = 1
			}

			if _, ok := a.ParticipantAnalyses[nameToFirstName(m.SenderName)].Mentions[word]; ok {
				a.ParticipantAnalyses[nameToFirstName(m.SenderName)].Mentions[word]++
			} else {
				a.ParticipantAnalyses[nameToFirstName(m.SenderName)].Mentions[word] = 1
			}
		}

		if _, ok := a.Words[word]; ok {
			a.Words[word]++
		} else {
			a.Words[word] = 1
		}

		if _, ok := a.ParticipantAnalyses[nameToFirstName(m.SenderName)].Words[word]; ok {
			a.ParticipantAnalyses[nameToFirstName(m.SenderName)].Words[word]++
		} else {
			a.ParticipantAnalyses[nameToFirstName(m.SenderName)].Words[word] = 1
		}
	}

	return nil
}

// StringFreq is a struct with a Value and its Frequency
type StringFreq struct {
	Value string
	Freq  int
}

// StringFreqs is for the array of StringFreq to be match the Sort interface
type StringFreqs []StringFreq

func (s StringFreqs) Len() int {
	return len(s)
}

func (s StringFreqs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s StringFreqs) Less(i, j int) bool {
	return s[i].Freq > s[j].Freq
}

// SortedAnalysis contains the aggregate analysis with the fields sorted
type SortedAnalysis struct {
	SortedParticipantAnalyses map[string]*SortedParticipantAnalysis
	Stickers                  StringFreqs
	Words                     StringFreqs
	Reactions                 StringFreqs
	Mentions                  StringFreqs
	MessageCount              int
}

// SortedParticipantAnalysis contains the sorted values of participants
type SortedParticipantAnalysis struct {
	Stickers     StringFreqs
	Words        StringFreqs
	Reactions    StringFreqs
	Mentions     StringFreqs
	MessageCount int
}

func newSortedAnalysis() SortedAnalysis {
	return SortedAnalysis{
		SortedParticipantAnalyses: make(map[string]*SortedParticipantAnalysis),
		Stickers:                  StringFreqs{},
		Words:                     StringFreqs{},
		Reactions:                 StringFreqs{},
		Mentions:                  StringFreqs{},
		MessageCount:              0,
	}
}

func newSortedParticipantAnalysis() *SortedParticipantAnalysis {
	return &SortedParticipantAnalysis{
		Stickers:     StringFreqs{},
		Words:        StringFreqs{},
		Reactions:    StringFreqs{},
		Mentions:     StringFreqs{},
		MessageCount: 0,
	}
}

// SortAnalysis returns the sorted analysis which each map sorted
func SortAnalysis(a Analysis) SortedAnalysis {
	s := newSortedAnalysis()
	s.Stickers = MapToSortedStringFreqs(a.Stickers)
	s.Words = MapToSortedStringFreqs(a.Words)
	s.Reactions = MapToSortedStringFreqs(a.Reactions)
	s.Mentions = MapToSortedStringFreqs(a.Mentions)
	s.MessageCount = a.MessageCount

	for k := range a.ParticipantAnalyses {
		s.SortedParticipantAnalyses[k] = newSortedParticipantAnalysis()
		s.SortedParticipantAnalyses[k].Stickers = MapToSortedStringFreqs(a.ParticipantAnalyses[k].Stickers)
		s.SortedParticipantAnalyses[k].Words = MapToSortedStringFreqs(a.ParticipantAnalyses[k].Words)
		s.SortedParticipantAnalyses[k].Reactions = MapToSortedStringFreqs(a.ParticipantAnalyses[k].Reactions)
		s.SortedParticipantAnalyses[k].Mentions = MapToSortedStringFreqs(a.ParticipantAnalyses[k].Mentions)
		s.SortedParticipantAnalyses[k].MessageCount = a.ParticipantAnalyses[k].MessageCount
	}

	return s
}

// MapToSortedStringFreqs generates and sorts the map to StringFreqs
func MapToSortedStringFreqs(m map[string]int) StringFreqs {
	sfs := StringFreqs{}

	for k, v := range m {
		sfs = append(sfs, StringFreq{
			Value: k,
			Freq:  v,
		})
	}

	sort.Sort(sfs)
	return sfs
}
