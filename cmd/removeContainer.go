/*
Copyright © 2024 Harsh Upadhyay amanupadhyay2004@gmail.com
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
)

var forceContainer bool

// removeContainerCmd represents the removeContainer command
var removeContainerCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove containers",
	Run: func(cmd *cobra.Command, args []string) {
		if forceContainer {
			forceRemoveContainer(args)
		} else {
			removeContainer(args)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeContainerCmd)
	removeContainerCmd.Flags().BoolVarP(&forceContainer, "force", "f", false, "Force removal of the container")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeContainerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeContainerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func removeContainer(ids []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Error connecting to Docker")
	}

	for _, id := range ids {	
		err := cli.ContainerRemove(ctx, id, containertypes.RemoveOptions{Force: false})
		if err != nil {
			fmt.Println("Container must be forced to remove")
			continue
		}
		fmt.Println("Container removed: " + id)
	}
}

func forceRemoveContainer(ids []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Error connecting to Docker")
	}

	for _, id := range ids {	
		err := cli.ContainerRemove(ctx, id, containertypes.RemoveOptions{Force: true})
		if err != nil {
			fmt.Println("Container must be forced to remove")
			continue
		}
		fmt.Println("Container removed: " + id)
	}
}