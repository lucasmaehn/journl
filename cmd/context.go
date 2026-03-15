/*
Copyright © 2026 LUCAS MÄHN <lucasmaehn@gmail.com>
*/
package cmd

import (
	"fmt"
	"regexp"

	"github.com/lucasmaehn/journl/config"
	"github.com/spf13/cobra"
)

var description string

// contextCmd represents the context command
var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage journl contexts",
	Long:  `Manage contexts to categorize your short entries (e.g., Work, Personal, Ideas).`,
}

var getContextCmd = &cobra.Command{
	Use:   "get",
	Short: "Show current active context",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Current context: %s\n", config.Get().ActiveContext)
	},
}

var setContextCmd = &cobra.Command{
	Use:   "set [context-name]",
	Short: "Sets the active context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		cur := config.Get()
		if _, set := cur.Contexts[name]; set {
			fmt.Printf("Set active context to %s\n", name)
			cur.ActiveContext = name
			if err := config.SetConfig(cur); err != nil {
				fmt.Printf("An error occure while updating active context: %v\n", err)
			}
		} else {
			fmt.Printf("Context with name %q does not exist\n", name)
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
		if description == "" {
			fmt.Printf("Adding context: %s\n", name)
		} else {
			fmt.Printf("Adding context: %s with description: %s\n", name, description)
		}

		cur := config.Get()
		context := config.Context{
			Name:        name,
			Description: description,
		}

		if _, set := cur.Contexts[name]; set {
			fmt.Printf("Context with name %q already exists. Context will not be added\n", name)
			return
		}

		cur.Contexts[name] = context

		if err := config.SetConfig(cur); err != nil {
			fmt.Printf("Failed to update config: %v\n", err)
		}
	},
}

func init() {
	contextCmd.AddCommand(getContextCmd)
	contextCmd.AddCommand(addContextCmd)
	contextCmd.AddCommand(setContextCmd)
	addContextCmd.Flags().StringVarP(&description, "description", "d", "", "Optional description of the context")

	rootCmd.AddCommand(contextCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// contextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// contextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
