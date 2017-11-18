package words

import (
	"bufio"
	"os"
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

/*
func Test_LDA(t *testing.T) {
	lda := NewSimpleLDA(&Configuration{})
	corpus := NewCorpus()

	processor := NewDefaultProcessor()
	processor.AddStrings(corpus, docs)

	err := lda.Init(corpus, 2, 0.1, 0.05)
	assert.Nil(t, err)

	_, err = lda.Train(1000)
	assert.Nil(t, err)
	lda.PrintTopWords(10)

	fmt.Printf("\n%+v", lda.topics)
}

*/

func Test_LDA(t *testing.T) {
	docs, err := readLines("./corpus/trump")
	assert.Nil(t, err)

	lda := NewSimpleLDA(&Configuration{})
	corpus := NewCorpus()

	processor := NewDefaultProcessor()
	processor.AddStrings(corpus, docs)

	err = lda.Init(corpus, 10, 0.1, 0.05)
	assert.Nil(t, err)

	_, err = lda.Train(10000)
	assert.Nil(t, err)
	lda.PrintTopWords(10)
}

func readLines(path string) ([]string, error) {
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
