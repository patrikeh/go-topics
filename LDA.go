package topics

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/golang/glog"
)

const (
	defaultAlphaSum = 50
	defaultBeta     = 0.01
)

// Implements a simple non-parallel Latent Dirichlet Allocation
type LDA struct {
	config  *Configuration
	rng     *rand.Rand
	topics  *Topics
	corpus  *Corpus
	alpha   float64 // Dir(alpha) - smoothing factor doc-topic distribution
	beta    float64 // Dir(beta) - smoothing factor topic-word distribution
	betaSum float64 // beta * size(vocabulary) constant
}

func NewLDA(config *Configuration) *LDA {
	rng := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	return &LDA{
		rng:    rng,
		config: config,
	}
}

func (l *LDA) Init(corpus *Corpus,
	numTopics int,
	alpha, beta float64) error {
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
		return fmt.Errorf("error initiating LDA - %s", err.Error())
	}

	return nil
}

// Initiate variables, MCMC set to random state
func (l *LDA) init(numTopics int) error {
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

// Gibb's sampling
func (l *LDA) Train(n int) (*Topics, error) {
	if l.topics == nil || l.corpus == nil {
		return nil, fmt.Errorf("unable to run LDA - uninitiated")
	}
	for it := 0; it < n; it++ {
		l.sample()

		if l.config.Verbose && it%l.config.PrintInterval == 0 && it != 0 {
			fmt.Printf("Iteration %d:\n", it)
			l.PrintTopWords(l.config.PrintNumWords)
		}
	}
	return l.topics, nil
}

// sample for all documents
func (l *LDA) sample() {
	for i, doc := range l.corpus.Documents {
		l.sampleDoc(i, doc)
	}
}

// sample per document
func (l *LDA) sampleDoc(di int, doc Document) {

	topicScores := make([]float64, l.topics.NumTopics, l.topics.NumTopics)
	for wi, word := range doc.Words {
		currTopic := l.topics.Topics[di][wi]
		wordTopic := l.topics.WordTopics[word]

		l.topics.DocTopics[di][currTopic]--
		l.topics.WordsPerTopic[currTopic]--
		wordTopic[currTopic]--

		var sum float64
		for topic := 0; topic < l.topics.NumTopics; topic++ {
			topicScores[topic] = (l.alpha + float64(l.topics.DocTopics[di][topic])) *
				((l.beta + float64(wordTopic[topic])) /
					(l.betaSum + float64(l.topics.WordsPerTopic[topic])))
			sum += topicScores[topic]
		}

		newTopic, err := l.sampleTopic(sum, topicScores)
		if err != nil {
			glog.Errorf("unable to sample topic for w(%d, %d) - %s", di, wi, err.Error())
		}

		l.assign(di, wi, newTopic)
	}
}

func (l *LDA) sampleTopic(sum float64, multinomial []float64) (int, error) {
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

func (l *LDA) assign(di, wi, topic int) {
	l.topics.Topics[di][wi] = topic
	l.topics.DocTopics[di][topic]++
	l.topics.WordTopics[l.corpus.Documents[di].Words[wi]][topic]++
	l.topics.WordsPerTopic[topic]++
}

type TopicWord struct {
	Occurrences int
	Word        int
}
type TopicWords []TopicWord

func (t TopicWords) Len() int           { return len(t) }
func (t TopicWords) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TopicWords) Less(i, j int) bool { return t[i].Occurrences < t[j].Occurrences }

func (l LDA) PrintTopWords(n int) {
	writer := tabwriter.NewWriter(os.Stdout, 8, 8, 2, ' ', 0)

	topicWords := make(TopicWords, l.topics.NumTokens, l.topics.NumTokens)
	fmt.Fprintln(writer, "Topic\tTokens\tWords")
	for topic := 0; topic < l.topics.NumTopics; topic++ {
		for w := 0; w < l.topics.NumTypes; w++ {
			topicWords[w] = TopicWord{
				Word:        w,
				Occurrences: l.topics.WordTopics[w][topic],
			}
		}
		sort.Sort(sort.Reverse(topicWords))

		n = min(n, l.topics.NumTypes)
		words := ""
		for i := 0; i < n; i++ {
			word := fmt.Sprintf("%s(%d)", l.corpus.Vocabulary.Words[topicWords[i].Word], topicWords[i].Occurrences)
			if i == 0 {
				words = word
			} else {
				words = fmt.Sprintf("%s %s", words, word)
			}
		}
		fmt.Fprintf(writer, "%d\t%d\t%s\n", topic, l.topics.WordsPerTopic[topic], words)
	}
	writer.Flush()
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
