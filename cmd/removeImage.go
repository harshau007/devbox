/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	// "strings"

	imagetype "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// removeImageCmd represents the removeImage command
var removeImageCmd = &cobra.Command{
	Use:   "rmi",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		removeImages(args)
	},
}

func init() {
	rootCmd.AddCommand(removeImageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeImageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeImageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func removeImages(ids []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Error connecting to Docker")
	}


	for _, id := range ids {	
		images, err := cli.ImageRemove(ctx, id, imagetype.RemoveOptions{Force: true, PruneChildren: true})
		if err != nil {
			fmt.Println("Image must be forced to remove")
		}
		fmt.Println("Deleted: " + images[1].Deleted[7:])
	}
}