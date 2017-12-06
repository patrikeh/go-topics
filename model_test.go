package topics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Vocabulary(t *testing.T) {
	v := NewVocabulary()

	v.Set("t")
	assert.Equal(t, v.Indices["t"], 0)
	assert.Equal(t, "t", v.Words[v.Indices["t"]])

	// Add twice does nothing
	v.Set("t")
	assert.Len(t, v.Words, 1)
}
