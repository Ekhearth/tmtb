package main

import (
    "fmt"
    "strings"
    "encoding/json"
    "io/ioutil"

    "mvdan.cc/xurls"
    "github.com/mb-14/gomarkov"
)

type trumpTweet struct {
  Text string `json:"text"`
}

func main() {
  // num is only for us to visually see the number of tweets processed in the console
  num := 1

  // Create our Markov chain
  chain := gomarkov.NewChain(1)

  //data, err := ioutil.ReadFile("trumpsample.json")
  data, err := ioutil.ReadFile("TrumpTweets.json")
  if err != nil {
		fmt.Println("ERROR reading from trumpsample.json")
	}

  s := []trumpTweet{}
  err = json.Unmarshal(data, &s)
  if err != nil {
    panic(err)
  }

  // This iterates over all the tweets and is our "foreach" statement
  for _, tweet := range s {
    fmt.Println(num)
    thisTweet := tweet.Text
    // gather any URLs in the current tweet
    urls := xurls.Strict().FindAllString(thisTweet, -1)

    // If a URL is found remove it. I'll need to think about the case where there are multiple URLs in a tweet.
    if len(urls) > 0 {
      for i := 0; i < len(urls); i++ {
        tweetSlice := strings.Fields(thisTweet)
        tweetSlice = remove(tweetSlice, urls[i])
        thisTweet = strings.Join(tweetSlice, " ")
      }
    }

    // Add this string to our Markov chain
    if thisTweet != "" {
      chain.Add(strings.Split(thisTweet, " "))
    }
    num++
  }

  // Generate our new tweet here - taken directly from the HNStory generator
  tokens := []string{gomarkov.StartToken}
  fmt.Printf("\n")
  for tokens[len(tokens)-1] != gomarkov.EndToken {
    //fmt.Printf("Token Length: %d\n", len(tokens))
		next, _ := chain.Generate(tokens[(len(tokens) - 1):])
		tokens = append(tokens, next)
	}
	fmt.Println(strings.Join(tokens[1:len(tokens)-1], " "))
}

// Our nice function to remove a word from a string
func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
