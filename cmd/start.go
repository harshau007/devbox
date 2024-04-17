/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		nameFlag, _ := cmd.Flags().GetString("name")
		if nameFlag == "" {
			fmt.Print("\nEnter the name of the container: ")
			fmt.Scanln(&nameFlag)
		}
		startContainer(nameFlag)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringP("name", "n", "","Name of the container")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func startContainer(name string) error {
	cmd := exec.Command("./start", name)
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Error executing the script: %v\n", err)
        return err
    }

    outputLines := strings.Split(strings.TrimSpace(string(output)), "\n")
    containerId := outputLines[len(outputLines)-1]
    fmt.Printf("Container created with ID: %s\n", containerId)

    return nil
}
