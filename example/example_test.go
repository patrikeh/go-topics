package example

import (
	"testing"

	topics "github.com/patrikeh/go-topics"
	"github.com/stretchr/testify/assert"
)

func Test_Example(t *testing.T) {
	processor := topics.NewProcessor(
		topics.Transformations{topics.ToLower, topics.RemoveTwitterUsernames, topics.Sanitize, topics.MinLen, topics.GetStopwordFilter("../stopwords/en"), topics.GetStopwordFilter("../stopwords/se")},
	)
	corpus, err := processor.ImportSingleFileCorpus(topics.NewCorpus(), "./corpus")
	assert.Nil(t, err)

	lda := topics.NewLDA(&topics.Configuration{Verbose: true, PrintInterval: 500, PrintNumWords: 8})
	err = lda.Init(corpus, 8, 0, 0)
	assert.Nil(t, err)

	_, err = lda.Train(10000)
	assert.Nil(t, err)
	lda.PrintTopWords(8)
}
