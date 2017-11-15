package words

type Inferer interface {
	Train(corpus Corpus, numIterations, numTopics int) *Topics
}

type Configuration struct {
	Parallellism int
	Debug        bool
}
