package words

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var docs = []string{
	"I like to eat broccoli and bananas yum.",
	"I ate a banana and spinach smoothie for breakfast.",
	"Chinchillas and kittens are cute.",
	"My sister adopted cute kittens yesterday.",
	"Look at this cute hamster munching on a piece of chinchillas.",
}

func Test_LDA(t *testing.T) {
	lda := NewSimpleLDA(&Configuration{})
	corpus := NewCorpus()

	processor := NewDefaultProcessor()
	processor.AddStrings(corpus, docs)

	err := lda.Init(corpus, 2, 0, 0)
	assert.Nil(t, err)

	_, err = lda.Train(1000)
	assert.Nil(t, err)
	lda.PrintTopWords(10)

	fmt.Printf("\n%+v", lda.topics)

}
