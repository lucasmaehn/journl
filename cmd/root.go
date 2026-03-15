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

	"github.com/lucasmaehn/journl/config"
	"github.com/lucasmaehn/journl/editor"
	"github.com/lucasmaehn/journl/logstore"
	"github.com/spf13/cobra"
)

var cfgFile string

var dbPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "journl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		stdin := getStdin()

		store, err := logstore.NewSQLite(config.Get().DBPath)
		if err != nil {
			return fmt.Errorf("failed to initialize logstore: %w", err)
		}

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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.journl/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.InitConfig(cfgFile)
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
