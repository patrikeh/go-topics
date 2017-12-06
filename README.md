# go-topics
A very basic LDA implementation with some convenient utilities.

## Usage
Create a processor from a set of transformations of the form ```func(word string) (new string, keep bool)```, that may filter or transform words:
```
processor := topics.NewProcessor(
  topics.Transformations{
    topics.ToLower, 
    topics.Sanitize, 
    topics.MinLen, 
    topics.GetStopwordFilter("../stopwords/en")})
```
Read data and apply transformations to build a corpus:
```
corpus, err := processor.ImportSingleFileCorpus(topics.NewCorpus(), "./corpus")
```
Run LDA and print the results:
```
lda := topics.NewLDA(&topics.Configuration{Verbose: true, PrintInterval: 500, PrintNumWords: 8})
err = lda.Init(corpus, 2, 0, 0) // Iterations, α, β hyperparameters

_, err = lda.Train(1000)
lda.PrintTopWords(5)
```
Resulting in something like:
```
Topic   Tokens  Words
0       9       like(1) eat(1) broccoli(1) bananas(1) ate(1)
1       14      cute(3) kittens(2) chinchillas(2) piece(1) look(1)
```
