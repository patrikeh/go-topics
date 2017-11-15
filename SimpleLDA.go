package words

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/golang/glog"
)

const (
	defaultAlphaSum = 50
	defaultBeta     = 0.01
)

// Implements a simple non-parallel Latent Dirichlet Allocation
type SimpleLDA struct {
	config  *Configuration
	rng     *rand.Rand
	topics  *Topics
	corpus  *Corpus
	alpha   float64 // Dir(alpha) - smoothing factor doc-topic distribution
	beta    float64 // Dir(beta) - smoothing factor topic-word distribution
	betaSum float64 // beta * size(vocabulary) constant
}

func NewSimpleLDA(config *Configuration) *SimpleLDA {
	rng := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	return &SimpleLDA{
		rng:    rng,
		config: config,
	}
}

func (l *SimpleLDA) Train(corpus *Corpus,
	numIterations, numTopics int,
	alpha, beta float64) (*Topics, error) {
	l.corpus = corpus
	if alpha == 0.0 {
		alpha = defaultAlphaSum / float64(numTopics)
	}
	if beta == 0.0 {
		beta = defaultBeta
	}
	l.alpha, l.beta, l.betaSum = alpha, beta, float64(corpus.Vocabulary.Size())*beta
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

	localTopics := make([]int, l.topics.NumTopics, l.topics.NumTopics)
	for wi := range doc.Words {
		localTopics[l.topics.Topics[di][wi]]++
	}

	topicScores := make([]float64, l.topics.NumTopics, l.topics.NumTopics)
	for wi, word := range doc.Words {
		currTopic := l.topics.Topics[di][wi]
		wordTopic := l.topics.WordTopics[word]

		localTopics[currTopic]--
		wordTopic[currTopic]--
		l.topics.WordsPerTopic[currTopic]--

		var sum float64
		for topic := 0; topic < l.topics.NumTopics; topic++ {
			topicScores[topic] = (l.alpha + float64(localTopics[topic])) *
				((l.beta+float64(wordTopic[topic]))/
					l.betaSum + float64(l.topics.WordsPerTopic[topic]))
			sum += topicScores[topic]
		}

		newTopic, err := l.sampleTopic(sum, topicScores)
		if err != nil {
			glog.Errorf("unable to sample topic for w(%d, %d) - %s", di, wi, err.Error())
		}

		l.assign(di, wi, newTopic)
	}
}

func (l *SimpleLDA) sampleTopic(sum float64, multinomial []float64) (int, error) {
	sample := l.rng.Float64() * sum
	n := len(multinomial)
	for sample > 0.0 {
		n--
		sample -= multinomial[n]
	}
	if n < 0 {
		return -1, fmt.Errorf("unable to sample topic")
	}
	return n, nil
}

func (l *SimpleLDA) assign(di, wi, topic int) {
	l.topics.Topics[di][wi] = topic
	l.topics.DocTopics[di][topic]++
	l.topics.WordTopics[l.corpus.Documents[di].Words[wi]][topic]++
	l.topics.WordsPerTopic[topic]++
}
