package words

type Corpus struct {
	Documents  []Document
	Vocabulary *Vocabulary
}

func NewCorpus() Corpus {
	return Corpus{
		Vocabulary: NewVocabulary(),
	}
}

type Document struct {
	Name  string
	Words []int // Vocabulary word indices
}

type Vocabulary struct {
	Words   []string
	Indices map[string]int
}

func NewVocabulary() *Vocabulary {
	return &Vocabulary{
		Words:   make([]string, 0),
		Indices: map[string]int{},
	}
}

func (v *Vocabulary) Set(w string) {
	if _, found := v.Indices[w]; found {
		return
	}
	v.Words = append(v.Words, w)
	v.Indices[w] = len(v.Words) - 1
}

type Topics struct {
	NumTypes, NumTokens, NumTopics int
	Topics                         [][]int // idx by <document, seq>
	DocTopics                      [][]int // idx by <document, topic>
	WordTopics                     [][]int // idx by <word, topic>
	WordsPerTopic                  []int   // idx by topic
}

func NewTopics(numTopics, numDocs, numTypes int) *Topics {
	return &Topics{
		NumTopics:     numTopics,
		NumTypes:      numTypes,
		Topics:        make([][]int, numDocs, numDocs),
		DocTopics:     initMatrix(numDocs, numTopics),
		WordTopics:    initMatrix(numTypes, numTopics),
		WordsPerTopic: make([]int, numTopics, numTopics),
	}
}

func initMatrix(m, n int) [][]int {
	matrix := make([][]int, m, m)
	for i := range matrix {
		matrix[i] = make([]int, n, n)
	}
	return matrix
}
