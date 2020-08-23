package main

import (
	"encoding/json"
	"flag"
	"log"
	"strconv"
	"sync"

	"golang.org/x/oauth2/clientcredentials"
	// other imports
	"fmt"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func keys(keyFile string) (key, secret string, err error) {
	var keys struct {
		Key    string `json:"consumer_key"`
		Secret string `json:"consumer_secret"`
	}
	f, err := os.Open(keyFile)
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.Decode(&keys)
	return keys.Key, keys.Secret, nil
}

func getClient(key string, secret string) *twitter.Client {

	if key == "" || secret == "" {
		log.Fatal("Application Access Token required")
	}
	// Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
	config := &clientcredentials.Config{
		ClientID:     key,
		ClientSecret: secret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}

	httpClient := config.Client(oauth2.NoContext)
	client := twitter.NewClient(httpClient)

	return client
}

func main() {
	fmt.Println("Go-Twitter Bot v0.01")
	var (
		keyFile string
	)
	flag.StringVar(&keyFile, "key", ".keys.json", "The file where you store your consumer key and secret for the Twitter API.")
	key, secret, err := keys(keyFile)
	if err != nil {
		panic(err)
	}

	client := getClient(key, secret)

	router := gin.Default()
	router.GET("/twitter/retweets/:user_handle/max", func(c *gin.Context) {
		userhandle := c.Param("user_handle")
		// fmt.Println(userhandle)
		tweetIDs := getTweets(client, userhandle)
		//fmt.Println(tweetIDs)
		retweetmap := mapretweets(client, tweetIDs)
		//fmt.Println(retweetmap)
		winners, max := findwinner(retweetmap)
		for k := range retweetmap {
			delete(retweetmap, k)
		}
		c.JSON(200, gin.H{
			"winnersofcontest": winners,
			"countofretweets":  max,
		})
	})

	router.GET("/twitter/tweet/:user_handle/latest", func(c *gin.Context) {
		userhandle := c.Param("user_handle")
		userTimelineParams := &twitter.UserTimelineParams{ScreenName: userhandle, ExcludeReplies: twitter.Bool(true), Count: 10}
		tweets, _, _ := client.Timelines.UserTimeline(userTimelineParams)
		for _, tweetinfo := range tweets {
			c.JSON(200, gin.H{
				"Tweet": tweetinfo.Text,
			})
		}

	})
	router.Run()
}

func getTweets(client *twitter.Client, userhandle string) []string {
	userTimelineParams := &twitter.UserTimelineParams{ScreenName: userhandle, ExcludeReplies: twitter.Bool(true), IncludeRetweets: twitter.Bool(false), Count: 100}
	tweets, _, _ := client.Timelines.UserTimeline(userTimelineParams)

	usertweets := make([]string, 0, len(tweets))
	for _, tweet := range tweets {
		usertweets = append(usertweets, tweet.IDStr)
	}
	return usertweets
}

func mapretweets(client *twitter.Client, tweetIDs []string) map[string]int {
	count := make(map[string]int)
	var wg sync.WaitGroup
	for i, tweet := range tweetIDs {
		wg.Add(1)
		temp, _ := strconv.ParseInt(tweet, 10, 64)
		go worker(i, temp, client, &count, &wg)
	}
	wg.Wait()

	return count
}

func worker(id int, twitID int64, client *twitter.Client, count *map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	twits, _, _ := client.Statuses.Retweets(twitID, &twitter.StatusRetweetsParams{Count: 100})
	for _, twit := range twits {
		(*count)[twit.User.Name] = (*count)[twit.User.Name] + 1
	}
}

func findwinner(retweetmap map[string]int) ([]string, int) {
	var maxretweets = 0
	var winnernames []string
	for key, value := range retweetmap {
		if maxretweets < value {
			winnernames = nil
			maxretweets = value
			winnernames = append(winnernames, key)
		} else if maxretweets == value {
			winnernames = append(winnernames, key)
		}
	}
	return winnernames, maxretweets
}
