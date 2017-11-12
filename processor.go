package words

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

type Filter func(string) bool
type Transformation func(string) string
type Filters []Filter
type Transformations []Transformation // i.e. stemming, lemmatization, etc
type Processor struct {
	filters         Filters
	transformations Transformations
}

func NewDefaultProcessor() *Processor {
	return NewProcessor(
		Filters{},
		Transformations{ToLower, Sanitize},
	)
}
func NewProcessor(filters Filters,
	transformations Transformations) *Processor {

	return &Processor{
		filters:         filters,
		transformations: transformations,
	}
}

func (p *Processor) AddDocs(corpus *Corpus, docs []io.Reader) *Corpus {

	for _, doc := range docs {
		err := p.AddDoc(corpus, doc)
		if err != nil {
			panic(err)
		}
	}

	return corpus
}

func (p *Processor) AddDoc(corpus *Corpus, doc io.Reader) error {
	document, err := p.processDocument(doc, "", corpus.Vocabulary)
	if err != nil {
		panic(err)
	}
	corpus.Documents = append(corpus.Documents, *document)
	return nil
}

func (p *Processor) processDocument(input io.Reader, name string, vocab *Vocabulary) (*Document, error) {
	s := bufio.NewScanner(input)
	s.Split(bufio.ScanWords)

	document := NewDocument(name)

	for s.Scan() {
		w := s.Text()
		w, ok := p.processWord(w)
		if !ok {
			continue
		}

		document.Add(vocab.Set(w))
	}
	return document, nil
}

func (p *Processor) processWord(w string) (string, bool) {
	if p.filter(w) {
		return "", false
	}
	return p.transform(w), true
}

func (p *Processor) filter(w string) bool {
	for _, f := range p.filters {
		if f(w) {
			return true
		}
	}
	return false
}

func (p *Processor) transform(w string) string {
	for _, t := range p.transformations {
		w = t(w)
	}
	return w
}

var StopWord = map[string]bool{
	"this": true,
	"I":    true,
	"that": true,
}

func StopWordFilter(word string) bool {
	if StopWord[word] {
		return false
	}
	return true
}

func ToLower(w string) string {
	return strings.ToLower(w)
}

var reg = regexp.MustCompile("[a-zA-Z]+")

func Sanitize(w string) string {
	return reg.FindString(w)
}
