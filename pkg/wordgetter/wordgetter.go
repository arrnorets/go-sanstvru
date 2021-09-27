package wordgetter

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Sanstv struct {
	guessedWord     string // A random word from sanstv.ru for our game
	wordDescription string //URL with guessedWord description
}

func (s *Sanstv) initialize() {

	url := "https://sanstv.ru/randomWord/lang-ru/strong-2/count-1/word-???????ajax=#result&lang=ru&strong=2&count=1&word="

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Get(url)
	// handle the error if there is one
	if err != nil {
		s.guessedWord = ""
		s.wordDescription = ""
		return
	}

	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		s.guessedWord = ""
		s.wordDescription = ""
		return
	}
	foundResult := false
	randomWordURI := ""
	var f func(*html.Node, bool, *string)
	f = func(n *html.Node, foundResult bool, word *string) {
		if n.Type == html.ElementNode && n.Data == "td" {
			for _, tag := range n.Attr {
				if tag.Val == "result" {
					foundResult = true
				}
			}
		}
		if n.Type == html.ElementNode && n.Data == "a" && foundResult {
			for _, tag := range n.Attr {
				if tag.Key == "href" {
					*word = tag.Val
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, foundResult, word)
		}
	}
	f(doc, foundResult, &randomWordURI)
	lastInd := strings.LastIndex(randomWordURI, "/")
	s.guessedWord = randomWordURI[lastInd+1:]
	s.wordDescription = "https://sanstv.ru" + randomWordURI
}

func (s *Sanstv) getWord() string {
	return s.guessedWord
}

func (s *Sanstv) getDesc() string {
	return s.wordDescription
}
