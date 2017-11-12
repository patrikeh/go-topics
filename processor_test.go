package words

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_processDocument(t *testing.T) {
	p := NewDefaultProcessor()

	cases := []struct {
		in                string
		expectedVocabSize int
		expectedSequence  []int
	}{
		{"this is a sentence", 4, []int{0, 1, 2, 3}},
		{"Is this a sentence", 4, []int{1, 0, 2, 3}},
		{"another sentence", 5, []int{4, 3}},
	}
	c := NewCorpus()

	for i, test := range cases {
		reader := strings.NewReader(test.in)
		assert.Nil(t, p.AddDoc(c, reader))
		assert.Len(t, c.Vocabulary.Words, test.expectedVocabSize)
		assert.Equal(t, test.expectedSequence, c.Documents[i].Words)
	}

	l := NewLDA(&Configuration{})
	topics, err := l.Train(c, 100, 10)
	assert.Nil(t, err)
	fmt.Printf("%+v", topics)
}
