/*
Copyright © 2026 LUCAS MÄHN <lucasmaehn@gmail.com>
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// clipCmd represents the clip command
var clipCmd = &cobra.Command{
	Use:   "clip",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// store, err := logstore.New()
		// if err != nil {
		// 	return err
		// }
		//
		// files, err := clip.Paste()
		// if err != nil {
		// 	return err
		// }
		//
		// var attachmentErrors []error
		// for _, fp := range files {
		// 	if err := createAttachment(fp); err != nil {
		// 		attachmentErrors = append(attachmentErrors, err)
		// 	}
		// }
		//
		return errors.New("clipboard functionality is not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(clipCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
