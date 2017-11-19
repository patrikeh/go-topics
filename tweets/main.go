package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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

	fileName := os.Args[1] + "/" + os.Args[2]

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, tweet := range tweets {
		fmt.Fprintln(w, tweet)
	}
	w.Flush()

	fmt.Printf("wrote %d tweets to %s\n", len(tweets), fileName)
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
	tweets := make([]string, 0, 0)

	n, rem := int(numTweets/100), numTweets%100

	var next int64
	for i := 0; i <= n; i++ {
		if i == n {
			if rem == 0 {
				break
			}
			numTweets = rem
		}
		retrieved, maxID, err := s.GetTweets(query, numTweets, next)
		next = maxID
		if err != nil {
			return nil, err
		}

		tweets = append(tweets, retrieved...)
		fmt.Printf("Got %d tweets\n", len(tweets))
	}

	return tweets, nil
}

func (s *Scraper) GetTweets(query string, numTweets int, maxID int64) ([]string, int64, error) {
	search, resp, err := s.twitterClient.Search.Tweets(&twitter.SearchTweetParams{
		Query: query,
		Count: numTweets,
		MaxID: maxID, // Equals lowest ID of already retrieved
		Lang:  "en",
	})
	defer resp.Body.Close()
	if err != nil {
		return nil, 0, err
	}

	tweets := make([]string, len(search.Statuses), len(search.Statuses))
	for i, tweet := range search.Statuses {
		tweets[i] = tweet.Text
	}

	return tweets, getNext(search.Metadata.NextResults), nil
}

var nextResultRegex = regexp.MustCompile(`\d+`)

func getNext(nextResult string) int64 {
	idString := nextResultRegex.FindString(nextResult)
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		panic(err)
	}
	return id
}
