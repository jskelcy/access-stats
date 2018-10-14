package types

import (
	"strings"
	"sync"
)

// Block represents logs parsed in 10 section time window.
// Data is indexed in Data by Section.
type Block struct {
	sync.Mutex
	Total int

	//histograms
	HistSection *Histogram
	Hist2XX     *Histogram
	Hist3XX     *Histogram
	Hist4XX     *Histogram
	Hist5XX     *Histogram
	HistUser    *Histogram
}

// NewBlock returns Block struct
func NewBlock() *Block {
	return &Block{
		HistSection: NewHistogram(),
		Hist2XX:     NewHistogram(),
		Hist3XX:     NewHistogram(),
		Hist4XX:     NewHistogram(),
		Hist5XX:     NewHistogram(),
		HistUser:    NewHistogram(),
	}
}

// Ingest parses a log line and updates Block data with
// relevant section metrics. This function is thread safe.
func (b *Block) Ingest(logLine string) {
	section := getSection(logLine)
	statusCode := getStatusCode(logLine)
	user := getUser(logLine)

	b.HistSection.Add(section, 1)
	b.addStatusCode(section, statusCode)
	b.HistUser.Add(user, 1)

	b.Lock()
	b.Total++
	defer b.Unlock()
}

func (b *Block) addStatusCode(section, statusCode string) {
	switch statusCode[0] {
	case '2':
		b.Hist2XX.Add(section, 1)
	case '3':
		b.Hist3XX.Add(section, 1)
	case '4':
		b.Hist4XX.Add(section, 1)
	case '5':
		b.Hist5XX.Add(section, 1)
	}
}

func getSection(data string) string {
	tokens := strings.Split(data, " ")
	section := tokens[6]
	sectionTokens := strings.Split(section, "/")
	if len(sectionTokens) < 4 {
		return strings.Join(sectionTokens, "/")
	}
	return strings.Join(sectionTokens[:3], "/")
}

func getStatusCode(data string) string {
	return strings.Split(data, " ")[8]
}

func getUser(data string) string {
	return strings.Split(data, " ")[2]
}
