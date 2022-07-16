package cmd

import (
	"github.com/hookart/twitter-mentions/listener"
	"github.com/hookart/twitter-mentions/routes"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the webserver",
	Run: func(cmd *cobra.Command, args []string) {
		go listener.Listen()
		routes.Serve()
	},
}
