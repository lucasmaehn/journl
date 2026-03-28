/*
Copyright © 2026 LUCAS MÄHN <lucasmaehn@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lucasmaehn/journl/editor"
	"github.com/lucasmaehn/journl/logstore"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stdin := getStdin()

		ctx, err := app.Config.ActiveContext()
		cobra.CheckErr(err)

		store, err := logstore.New(ctx.Name, ctx.Store)
		cobra.CheckErr(err)

		var text string
		if len(args) == 0 {
			reader, err := editor.Open()
			if err != nil {
				return err
			}
			t, err := io.ReadAll(reader)
			if err != nil {
				return err
			}
			text = string(t)
		} else {
			text = strings.Join(args, " ")
		}

		if len(strings.TrimSpace(text)) == 0 {
			return errors.New("skipping log entry because message was empty")
		}

		fmt.Println("stdin length", len(stdin))
		if len(stdin) > 0 {
			text += "\n"
			text += stdin
		}

		return store.Commit(text)
	},
}

func getStdin() string {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		bytes, err := io.ReadAll(os.Stdin)
		if err == nil {
			return string(bytes)
		}
	}
	return ""
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
