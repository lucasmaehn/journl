/*
Copyright © 2026 LUCAS MÄHN <lucasmaehn@gmail.com>
*/
package cmd

import (
	"github.com/lucasmaehn/journl/config"
	"github.com/lucasmaehn/journl/logstore"
	"github.com/lucasmaehn/journl/ui"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ls, err := logstore.NewSQLite(config.Get().DBPath)
		if err != nil {
			return err
		}

		entries, err := ls.List()
		if err != nil {
			return err
		}

		ui.RenderJournal(entries)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
