package cmd

import (
	"log"

	twitterstream "github.com/fallenstedt/twitter-stream"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "configure twitter stream rules",
	Run: func(cmd *cobra.Command, args []string) {
		tok, err := twitterstream.NewTokenGenerator().SetApiKeyAndSecret(viper.GetString("twitter_client_id"), viper.GetString("twitter_client_secret")).RequestBearerToken()
		if err != nil {
			log.Fatalf(err.Error())
		}
		api := twitterstream.NewTwitterStream(tok.AccessToken)
		rules := twitterstream.NewRuleBuilder().
			AddRule("@HookProtocol \"gotrekt.xyz\"", "mentions protocol and gotrekt.xyz").Build()

		res, err := api.Rules.Create(rules, false)
		log.Println(res.Data)
	},
}
