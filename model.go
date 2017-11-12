package words

type Corpus struct {
	Documents  []Document
	Vocabulary *Vocabulary
}

func NewCorpus() *Corpus {
	return &Corpus{
		Documents:  []Document{},
		Vocabulary: NewVocabulary(),
	}
}

type Document struct {
	Name  string
	Words []int // sequence - vocabulary word indices
}

func NewDocument(name string) *Document {
	return &Document{
		Name:  name,
		Words: make([]int, 0, 0),
	}
}

func (d *Document) Add(word int) {
	d.Words = append(d.Words, word)
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

func (v *Vocabulary) Set(w string) int {
	if idx, found := v.Indices[w]; found {
		return idx
	}
	v.Words = append(v.Words, w)
	idx := len(v.Words) - 1
	v.Indices[w] = idx
	return idx
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
