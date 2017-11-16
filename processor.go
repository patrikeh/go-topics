package words

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type Transformation func(word string) (new string, keep bool)
type Transformations []Transformation // i.e. stemming, lemmatization, etc
type Processor struct {
	transformations Transformations
}

func NewDefaultProcessor() *Processor {
	return NewProcessor(
		Transformations{ToLower, Sanitize, PruneStopWord},
	)
}
func NewProcessor(transformations Transformations) *Processor {
	return &Processor{
		transformations: transformations,
	}
}

func (p *Processor) AddStrings(corpus *Corpus, docs []string) *Corpus {
	for _, doc := range docs {
		err := p.AddDoc(corpus, strings.NewReader(doc))
		if err != nil {
			panic(err)
		}
	}

	return corpus
}

func (p *Processor) AddReaders(corpus *Corpus, docs []io.Reader) *Corpus {
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
		w, keep := p.transform(w)
		if !keep {
			continue
		}
		document.Add(vocab.Set(w))
	}
	return document, nil
}

func (p *Processor) transform(w string) (string, bool) {
	var keep bool
	for _, t := range p.transformations {
		w, keep = t(w)
		if !keep {
			return "", false
		}
	}
	return w, true
}

var StopWord = func(path string) map[string]bool {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		fmt.Printf("%+v", err)
		return nil
	}

	s := bufio.NewScanner(f)

	words := make(map[string]bool, 0)
	for s.Scan() {
		words[s.Text()] = true
	}
	return words
}("stopwords/english")

func PruneStopWord(w string) (string, bool) {
	if StopWord[w] {
		return "", false
	}
	return w, true
}

func ToLower(w string) (string, bool) {
	return strings.ToLower(w), true
}

var reg = regexp.MustCompile("[a-zA-Z']+")

func Sanitize(w string) (string, bool) {
	return reg.FindString(w), true
}
