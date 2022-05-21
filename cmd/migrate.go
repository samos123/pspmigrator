package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Interactive command to migrate from PSP to PSA ",
	Long: `The interactive command will help with setting a suggested a
	Suggested Pod Security Standard for each namespace. In addition, it also
	checks whether a PSP object is mutating pods in every namespace.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Migrate")
	},
}
