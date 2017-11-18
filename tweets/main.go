package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	if len(os.Args) != 4 {
		println("tweets <output_path> <topic> <num_tweets>")
		return
	}

	scraper := NewScraper("zTQMkDtKiv4QlacGik0u8a3A5", os.Getenv("TWITTER_API_SECRET"),
		"338704968-yRuBFS5HI5RBsVwFmL9Myy25eybGzQas7GbVfPVD", os.Getenv("TWITTER_ACCESS_SECRET"))

	numTweets, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(err)
	}
	tweets, err := scraper.Collect(os.Args[2], numTweets)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(os.Args[1] + "/" + os.Args[2])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, tweet := range tweets {
		fmt.Fprintln(w, tweet)
	}
	w.Flush()
}

type Scraper struct {
	twitterClient *twitter.Client
}

func NewScraper(apiKey, apiSecret, accessToken, accessSecret string) *Scraper {
	config := oauth1.NewConfig(apiKey, apiSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return &Scraper{
		twitterClient: twitter.NewClient(httpClient),
	}
}

func (s *Scraper) Collect(query string, numTweets int) ([]string, error) {
	search, resp, err := s.twitterClient.Search.Tweets(&twitter.SearchTweetParams{
		Query:      query,
		Count:      numTweets,
		ResultType: "recent",
	})
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	tweets := make([]string, len(search.Statuses), len(search.Statuses))
	for i, tweet := range search.Statuses {
		tweets[i] = tweet.Text
	}
	return tweets, nil
}
