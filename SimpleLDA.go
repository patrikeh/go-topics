package words

import (
	"fmt"
	"math/rand"
	"time"
)

// Implements a simple non-parallel Latent Dirichlet Allocation
type SimpleLDA struct {
	config *Configuration
	rng    *rand.Rand
	topics *Topics
	corpus *Corpus
}

func NewSimpleLDA(config *Configuration) *SimpleLDA {
	rng := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	return &SimpleLDA{
		rng:    rng,
		config: config,
	}
}

func (l *SimpleLDA) Train(corpus *Corpus, numIterations, numTopics int) (*Topics, error) {
	l.corpus = corpus
	err := l.init(numTopics)
	if err != nil {
		return nil, fmt.Errorf("error initiating SimpleLDA - %s", err.Error())
	}

	// Gibb's sampling
	for it := 0; it < numIterations; it++ {
		l.sample()
	}

	return l.topics, nil
}

// Initiate variables, MCMC set to random state
func (l *SimpleLDA) init(numTopics int) error {
	if l.corpus == nil || l.corpus.Vocabulary == nil {
		return fmt.Errorf("missing corpus or vocabulary")
	}
	l.topics = NewTopics(numTopics, len(l.corpus.Documents), len(l.corpus.Vocabulary.Words))

	for di, doc := range l.corpus.Documents {
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

// sample for all documents
func (l *SimpleLDA) sample() {
	for i, doc := range l.corpus.Documents {
		l.sampleDoc(i, doc)
	}
}

// sample per document
func (l *SimpleLDA) sampleDoc(di int, doc Document) {
	var newTopic int

	for wi := range doc.Words {

		l.assign(di, wi, newTopic)
	}
}

func (l *SimpleLDA) assign(di, wi, topic int) {
	l.topics.Topics[di][wi] = topic
	l.topics.DocTopics[di][topic]++
	l.topics.WordTopics[l.corpus.Documents[di].Words[wi]][topic]++
	l.topics.WordsPerTopic[topic]++
}
