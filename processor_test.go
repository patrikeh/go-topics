package topics

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_processDocument(t *testing.T) {
	p := NewProcessor(Transformations{ToLower, Sanitize})
	cases := []struct {
		in                string
		expectedVocabSize int
		expectedSequence  []int
	}{
		{"this is a sentence.", 4, []int{0, 1, 2, 3}},
		{"Is this a sentence?", 4, []int{1, 0, 2, 3}},
		{"another, sentence", 5, []int{4, 3}},
	}
	c := NewCorpus()

	for i, test := range cases {
		reader := strings.NewReader(test.in)
		assert.Nil(t, p.AddDoc(c, reader))
		assert.Len(t, c.Vocabulary.Words, test.expectedVocabSize)
		assert.Equal(t, test.expectedSequence, c.Documents[i].Words)
	}
}
