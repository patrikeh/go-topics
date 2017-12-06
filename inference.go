package topics

type Inferer interface {
	Train(corpus Corpus, numIterations, numTopics int) *Topics
}

type Configuration struct {
	Parallelism   int
	Verbose       bool
	PrintInterval int
	PrintNumWords int
}
