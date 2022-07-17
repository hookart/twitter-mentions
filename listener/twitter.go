package listener

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	twitterstream "github.com/fallenstedt/twitter-stream"
	"github.com/fallenstedt/twitter-stream/stream"
	"github.com/hookart/twitter-mentions/models"
	"github.com/spf13/viper"
)

type StreamDataExample struct {
	Data struct {
		Text      string    `json:"text"`
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		AuthorID  string    `json:"author_id"`
	} `json:"data"`
	Includes struct {
		Users []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"users"`
	} `json:"includes"`
	MatchingRules []struct {
		ID  string `json:"id"`
		Tag string `json:"tag"`
	} `json:"matching_rules"`
}

func fetchTweets() stream.IStream {
	tok, err := twitterstream.NewTokenGenerator().SetApiKeyAndSecret(viper.GetString("twitter_client_id"), viper.GetString("twitter_client_secret")).RequestBearerToken()
	if err != nil {
		log.Fatalf(err.Error())
	}
	api := twitterstream.NewTwitterStream(tok.AccessToken).Stream

	api.SetUnmarshalHook(func(bytes []byte) (interface{}, error) {
		data := StreamDataExample{}

		if err := json.Unmarshal(bytes, &data); err != nil {
			fmt.Printf("failed to unmarshal bytes: %v", err)
		}

		return data, err
	})

	streamExpansions := twitterstream.NewStreamQueryParamsBuilder().
		AddExpansion("author_id").
		AddTweetField("created_at").
		Build()

	err = api.StartStream(streamExpansions)

	if err != nil {
		log.Fatalf("err")
	}
	return api
}
func Listen() {
	db := models.GetDBConnection()

	var re = regexp.MustCompile(`(?m)(0x[A-Za-z0-9]{40})`)

	fmt.Println("Starting Stream")

	// Start the stream
	// And return the library's api
	api := fetchTweets()

	// When the loop below ends, restart the stream
	defer Listen()

	// Start processing data from twitter after starting the stream
	for result := range api.GetMessages() {

		// Handle disconnections from twitter
		// https://developer.twitter.com/en/docs/twitter-api/tweets/volume-streams/integrate/handling-disconnections
		if result.Err != nil {
			fmt.Printf("got error from twitter: %v", result.Err)

			// Notice we "StopStream" and then "continue" the loop instead of breaking.
			// StopStream will close the long running GET request to Twitter's v2 Streaming endpoint by
			// closing the `GetMessages` channel. Once it's closed, it's safe to perform a new network request
			// with `StartStream`
			api.StopStream()
			continue
		}
		tweet := result.Data.(StreamDataExample)

		// fmt.Println(tweet.Data.Text)
		// bytes, _ := json.MarshalIndent(tweet, ">", "    ")
		// fmt.Println(string(bytes))

		for _, match := range re.FindAllString(tweet.Data.Text, -1) {
			log.Println("tweet matched!! ", tweet.Data.ID)
			username := ""
			for _, j := range tweet.Includes.Users {
				if j.ID == tweet.Data.AuthorID {
					username = j.Username
					log.Println("matched user name ", j.Name, j.Username)
				}
			}
			validation := models.Verification{}
			err := db.Find(&validation, &models.Verification{VerificationString: match}).Error
			if err == nil {
				account := models.Account{}
				db.Find(&account, "id = ?", validation.AccountID)

				account.TwitterHandle = username
				account.Verified = true
				log.Println("tweet updated!! ", account.ID)
				db.Save(&account)
				db.Delete(&validation)
			}
		}
	}
}
