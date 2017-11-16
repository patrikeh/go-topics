package words

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var docs = []string{
	"I like to eat broccoli and bananas.",
	"I ate a banana and spinach smoothie for breakfast.",
	"Chinchillas and kittens are cute.",
	"My sister adopted a kitten yesterday.",
	"Look at this cute hamster munching on a piece of broccoli.",
}

func Test_LDA(t *testing.T) {
	lda := NewSimpleLDA(&Configuration{})
	corpus := NewCorpus()

	processor := NewDefaultProcessor()
	processor.AddStrings(corpus, docs)

	err := lda.Init(corpus, 2, 0, 0)
	assert.Nil(t, err)

	_, err = lda.Train(4)
	assert.Nil(t, err)
	lda.PrintTopWords(10)

	fmt.Printf("\n%+v", lda.topics)

}
