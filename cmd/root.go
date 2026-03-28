/*
Copyright © 2026 LUCAS MÄHN <lucasmaehn@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucasmaehn/journl/config"
	"github.com/spf13/cobra"
)

var cfgFile string

var dbPath string

type App struct {
	Config *config.Config
}

var app *App

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "journl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		skipConfigLoad := map[string]bool{
			"init": true,
		}

		if skipConfigLoad[cmd.Name()] {
			return nil
		}

		cfg, err := config.Load(cfgFile)
		if err != nil {
			if errors.Is(err, config.ErrConfigNotFound) {
				return fmt.Errorf("Config file not found. Initialize a new, default config using `journl init` or pass the path to a valid config file using the `--config` option")
			}
			return err
		}

		app = &App{
			Config: cfg,
		}

		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	defaultCfg := filepath.Join(home, ".journl", "config.yaml")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", defaultCfg, "config file (default is $HOME/.journl/config.yaml)")
}
