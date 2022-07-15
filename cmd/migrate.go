package cmd

import (
	"github.com/hookart/twitter-mentions/models/migration"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "perform the migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migration.Migrate()
	},
}
