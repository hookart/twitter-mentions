package listener

import (
	"log"
	"regexp"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/hookart/twitter-mentions/models"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func Listen() {
	db := models.GetDBConnection()

	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     viper.GetString("twitter_client_id"),
		ClientSecret: viper.GetString("twitter_client_secret"),
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)

	// Twitter client
	client := twitter.NewClient(httpClient)

	params := &twitter.StreamFilterParams{
		Track:         []string{"@HookProtocol"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(params)

	if err != nil {
		log.Fatalf("err")
	}

	var re = regexp.MustCompile(`(?m)(0x[A-Za-z0-9]{40})`)
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		for _, match := range re.FindAllString(tweet.Text, -1) {
			validation := models.Verification{}
			err := db.Find(&validation, &models.Verification{VerificationString: match})
			if err == nil {
				account := models.Account{}
				db.Find(&account, "id = ?", validation.AccountID)

				account.TwitterHandle = tweet.User.IDStr
				account.Verified = true
				db.Save(&account)
				db.Delete(&validation)
			}
		}
	}

	for message := range stream.Messages {
		demux.Handle(message)
	}

}
