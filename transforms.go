package topics

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Transformation func(word string) (new string, keep bool)
type Transformations []Transformation // i.e. stemming, lemmatization, etc
type Processor struct {
	transformations Transformations
}

func GetStopwordFilter(path string) Transformation {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		fmt.Printf("%+v", err)
		return nil
	}

	s := bufio.NewScanner(f)

	stopWords := make(map[string]bool, 0)
	for s.Scan() {
		stopWords[s.Text()] = true
	}

	return func(stopWords map[string]bool) Transformation {
		return func(w string) (string, bool) {
			if stopWords[w] {
				return "", false
			}
			return w, true
		}
	}(stopWords)
}

func ToLower(w string) (string, bool) {
	return strings.ToLower(w), true
}

var wordReg = regexp.MustCompile("[a-zA-Z'åäöÅÄÖ]+")

func Sanitize(w string) (string, bool) {
	return strings.TrimSpace(wordReg.FindString(w)), true
}

func MinLen(w string) (string, bool) {
	if len(w) < 2 {
		return "", false
	}
	return w, true
}

func RemoveTwitterUsernames(w string) (string, bool) {
	if strings.Contains(w, "@") {
		return "", false
	}
	return w, true
}
