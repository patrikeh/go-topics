package words

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LDA(t *testing.T) {
	lda := NewLDA(&Configuration{})
	corpus := NewCorpus()
	_, err := lda.Train(corpus, 100, 10)
	assert.Nil(t, err)
}
