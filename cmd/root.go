/*
Copyright © 2024 NAME HERE amanupadhyay2004@gmail.com

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "devctl",
	Short: "DevControl is a powerful CLI tool to create and manage isolated containers for developers with their desired technology stacks, such as Node.js, Python, Rust, and more.",
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
	rootCmd.AddCommand(createcmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Version = "0.3.0-beta"
	createcmd.PersistentFlags().StringP("name", "n", "","Name of the container")
	createcmd.PersistentFlags().StringP("package", "p", "","Name of the package")
	createcmd.PersistentFlags().StringP("volume", "v", "","Path to the volume")
}


