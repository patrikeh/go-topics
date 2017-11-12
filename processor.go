package words

import (
	"regexp"
)

type Filter func(string) bool
type Transformation func(string) string
type Filters []Filter
type Transformations []Transformation // i.e. stemming, lemmatization, etc
type Processor struct {
	regex           *regexp.Regexp
	filters         Filters
	transformations Transformations
}

func NewProcessor(wordRegex string,
	filters Filters,
	transformations Transformations) *Processor {
	regex, _ := regexp.Compile(wordRegex)
	return &Processor{
		regex:           regex,
		filters:         filters,
		transformations: transformations,
	}
}

/*
func (p *Processor) FromFiles(documents []os.File) *Corpus {

}
func (p *Processor) FromStrings(documents []string) *Corpus {

}
*/
/*
FromFiles
FromStrings
*/
/* func (p *Processor) Process(documents []io.Reader) Corpus {
	corpus := NewCorpus()
	for _, w := range something {
		if p.filter(w) {
			continue
		}
		w = p.transform(w)
	}

}
*/
// Process filters & transforms input, builds vocab
/* func (p *Processor) ProcessFiles(files []os.File) Collection {

} */

func (p *Processor) filter(w string) bool {
	for _, f := range p.filters {
		if f(w) {
			return false
		}
	}
	return true
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
