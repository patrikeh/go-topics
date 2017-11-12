package words

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LDA(t *testing.T) {
	lda := NewLDA(&Configuration{})
	topics, err := lda.Train(&Corpus{}, 100, 10)
	assert.Nil(t, err)
	fmt.Printf("%+v", topics)
}
