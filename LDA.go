package words

import (
	"fmt"
	"math/rand"
	"time"
)

type LDA struct {
	config *Configuration
	rng    *rand.Rand
	topics *Topics
	corpus *Corpus
}

func NewLDA(config *Configuration) *LDA {
	rng := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	return &LDA{
		rng:    rng,
		config: config,
	}
}

func (l LDA) Train(corpus *Corpus, numIterations, numTopics int) (*Topics, error) {
	l.corpus = corpus
	err := l.init(corpus, numTopics)
	if err != nil {
		return nil, fmt.Errorf("error initiating LDA - %s", err.Error())
	}

	return l.topics, nil
}

// Initiate variables, MCMC set to random state
func (l LDA) init(corpus *Corpus, numTopics int) error {
	if corpus == nil || corpus.Vocabulary == nil {
		return fmt.Errorf("missing corpus or vocabulary")
	}
	l.topics = NewTopics(numTopics, len(corpus.Documents), len(corpus.Vocabulary.Words))

	for di, doc := range corpus.Documents {
		seqLen := len(doc.Words)
		l.topics.NumTokens += seqLen
		l.topics.Topics[di] = make([]int, seqLen, seqLen)

		for wi := 0; wi < seqLen; wi++ {
			topic := l.rng.Intn(l.topics.NumTopics)
			l.assign(di, wi, topic)
		}

	}
	return nil
}

func (l LDA) assign(di, wi, topic int) {
	l.topics.Topics[di][wi] = topic
	l.topics.DocTopics[di][topic]++
	l.topics.WordTopics[l.corpus.Documents[di].Words[wi]][topic]++
	l.topics.WordsPerTopic[topic]++
}
