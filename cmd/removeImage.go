/*
Copyright © 2024 Harsh Upadhyay amanupadhyay2004@gmail.com
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

var forceImage bool

// removeImageCmd represents the removeImage command
var removeImageCmd = &cobra.Command{
	Use:   "rmi",
	Short: "Remove images",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if forceImage {
			removeImagesForce(args)
		} else {
			removeImages(args)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeImageCmd)
	removeImageCmd.Flags().BoolVarP(&forceImage, "force", "f", false, "Force removal of the image")
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
		images, err := cli.ImageRemove(ctx, id, imagetype.RemoveOptions{Force: false, PruneChildren: false})
		if err != nil {
			fmt.Println("Image must be forced to remove")
			continue
		}
		if len(images[1].Deleted) > 0 {
			fmt.Println("No images were deleted")
		} else {
			fmt.Println("Deleted: " + images[1].Deleted[7:])
		}
	}
}

func removeImagesForce(ids []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Error connecting to Docker")
	}

	for _, id := range ids {
		images, err := cli.ImageRemove(ctx, id, imagetype.RemoveOptions{Force: true, PruneChildren: true})
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(images[1].Deleted) > 0 {
			fmt.Println("No images were deleted")
		} else {
			fmt.Println("Deleted: " + images[1].Deleted[7:])
		}
	}
}
