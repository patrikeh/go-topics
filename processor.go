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
		Transformations{ToLower, Sanitize, GetStopwordFilter("stopwords/en")},
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
		if !keep || w == "" {
			return "", false
		}
	}
	return w, true
}

// Process one document corpus
func (p *Processor) ImportSingleFileCorpus(corpus *Corpus, path string) (*Corpus, error) {
	docs, err := p.importSingleFileCorpus(path)
	if err != nil {
		return nil, err
	}

	return p.AddStrings(corpus, docs), err
}

// Read one document per line
func (p *Processor) importSingleFileCorpus(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func GetStopwordFilter(path string) Transformation {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		fmt.Printf("%+v", err)
		return nil
	}

	s := bufio.NewScanner(f)

	stopWords := make(map[string]bool, 0)
	for s.Scan() {
		stopWords[s.Text()] = true
	}

	return func(stopWords map[string]bool) Transformation {
		return func(w string) (string, bool) {
			if stopWords[w] {
				return "", false
			}
			return w, true
		}
	}(stopWords)
}

func ToLower(w string) (string, bool) {
	return strings.ToLower(w), true
}

var wordReg = regexp.MustCompile("[a-zA-Z'åäöÅÄÖ]+")

func Sanitize(w string) (string, bool) {
	return strings.TrimSpace(wordReg.FindString(w)), true
}

func MinLen(w string) (string, bool) {
	if len(w) < 2 {
		return "", false
	}
	return w, true
}

func RemoveTwitterUsernames(w string) (string, bool) {
	if strings.Contains(w, "@") {
		return "", false
	}
	return w, true
}
