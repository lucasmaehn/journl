/*
Copyright © 2026 LUCAS MÄHN <lucasmaehn@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"slices"

	"github.com/lucasmaehn/journl/config"
	"github.com/spf13/cobra"
)

var (
	description  string
	storeBackend string
)

// contextCmd represents the context command
var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage journl contexts",
	Long:  `Manage contexts to categorize your short entries (e.g., Work, Personal, Ideas).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Current context: %s\n", app.Config.CurrentContext)
	},
}

var useContextCmd = &cobra.Command{
	Use:   "use [context-name]",
	Short: "Changes the active context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		cfg := app.Config
		err := cfg.UseContext(name)
		cobra.CheckErr(err)

		err = config.Save(app.Config, cfgFile)
		cobra.CheckErr(err)
	},
}

var listContextCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all available contexts",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := app.Config

		contexts := cfg.ListContexts()

		names := make([]string, 0, len(contexts))

		for name := range contexts {
			names = append(names, name)
		}
		slices.Sort(names)

		for _, name := range names {
			if name == cfg.CurrentContext {
				fmt.Fprintf(os.Stdout, "* %s\n", name)
			} else {
				fmt.Fprintf(os.Stdout, "  %s\n", name)
			}
		}
	},
}

var addContextCmd = &cobra.Command{
	Use:   "add [context-name] --description [description]",
	Short: "Add a new context",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		name := args[0]
		validate := regexp.MustCompile(`^[a-z0-9_]+$`)
		if !validate.MatchString(name) {
			return fmt.Errorf("invalid context name '%s': must only contain lowercase letters and underscores", name)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		context := config.Context{
			Name:        name,
			Description: description,
		}

		switch config.StoreFormat(storeBackend) {
		case config.StoreFormatJSONL:
			context.Store = config.StoreConfig{
				Format: config.StoreFormatJSONL,
				Path:   path.Join(path.Dir(cfgFile), "journl.jsonl"),
				JSONL:  &config.JSONLConfig{},
			}
		case config.StoreFormatSQLite:
			context.Store = config.StoreConfig{
				Format: config.StoreFormatSQLite,
				Path:   path.Join(path.Dir(cfgFile), "journl.sqlite"),
				SQLite: &config.SQLiteConfig{},
			}
		default:
			cobra.CheckErr(fmt.Errorf("invalid store backend: %q", storeBackend))
		}

		cobra.CheckErr(app.Config.AddContext(name, context))
		cobra.CheckErr(config.Save(app.Config, cfgFile))
	},
}

func init() {
	contextCmd.AddCommand(addContextCmd)
	contextCmd.AddCommand(useContextCmd)
	contextCmd.AddCommand(listContextCmd)
	addContextCmd.Flags().StringVarP(&description, "description", "d", "", "Optional description of the context")
	addContextCmd.Flags().StringVarP(&storeBackend, "store", "s", string(config.StoreFormatSQLite), "Backend for storing journl entries. Further configuration is currently unavailable via CLI. Edit config file directly for this.")

	rootCmd.AddCommand(contextCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// contextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// contextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
