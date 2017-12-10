# go-topics
A very basic LDA (Latent Dirichlet Allocation) implementation with some convenient utilities.

## Usage
Create a processor from a set of transformations of the form ```func(word string) (new string, keep bool)```:
```go
processor := topics.NewProcessor(
  topics.Transformations{
    topics.ToLower, 
    topics.Sanitize, 
    topics.MinLen, 
    topics.GetStopwordFilter("../stopwords/en")})
```
Read data and apply transformations to build a corpus:
```go
var docs = []string{
	"I like to eat broccoli and bananas.",
	"I ate a banana and spinach smoothie for breakfast.",
	"Chinchillas and kittens are cute.",
	"My sister adopted cute kittens yesterday.",
	"Look at this cute hamster munching on a piece of chinchillas.",
}

corpus, err := processor.AddStrings(topics.NewCorpus(), docs)
```
Run LDA and print the results:
```go
lda := topics.NewLDA(&topics.Configuration{Verbose: true, PrintInterval: 500, PrintNumWords: 8})
err = lda.Init(corpus, 2, 0, 0) // K (number of topics), α, β (Dirichlet distribution smoothing factors)

_, err = lda.Train(1000) // Number of iterations
lda.PrintTopWords(5)
```
Resulting in something like:
```
Topic   Tokens  Words
0       9       like(1) eat(1) broccoli(1) bananas(1) ate(1)
1       14      cute(3) kittens(2) chinchillas(2) piece(1) look(1)
```
