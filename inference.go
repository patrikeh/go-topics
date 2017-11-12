package words

type Inferer interface {
	Train(corpus Corpus, numIterations, numTopics int) *Topics
}

type Configuration struct {
	Alpha        float64
	Beta         float64
	Parallellism int
}
