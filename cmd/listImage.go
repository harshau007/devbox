/*
Copyright © 2024 Harsh Upadhyay amanupadhyay2004@gmail.com
*/
package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/filters"
	imagetype "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// listImageCmd represents the listImage command
var listImageCmd = &cobra.Command{
	Use:   "images",
	Short: "List images",
	Run: func(cmd *cobra.Command, args []string) {
		listImages()
	},
}

func init() {
	rootCmd.AddCommand(listImageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listImageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listImageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listImages() error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()
	filters := filters.NewArgs(filters.Arg("label", "createdBy=DevBox"))
	images, err := cli.ImageList(ctx, imagetype.ListOptions{Filters: filters})
	if err != nil {
		return err
	}

	const (
		repositoryWidth = 30
		tagWidth        = 20
		imageIDWidth    = 20
		createdWidth    = 25
		sizeWidth       = 15
	)

	if len(images) == 0 {
		fmt.Println("No images found.")
		return nil
	}

	fmt.Println(strings.Repeat("-", repositoryWidth+tagWidth+imageIDWidth+createdWidth+sizeWidth+4))
	fmt.Printf("   %-*s %-*s %-*s %-*s %-*s\n",
		repositoryWidth, "REPOSITORY",
		tagWidth, "TAG",
		imageIDWidth, "IMAGE ID",
		createdWidth, "CREATED",
		sizeWidth, "SIZE")
	fmt.Println(strings.Repeat("-", repositoryWidth+tagWidth+imageIDWidth+createdWidth+sizeWidth+4))
	for _, image := range images {
		var imageName, tag string
		if len(image.RepoDigests) > 0 && len(image.RepoTags) > 0 {
			imageNameParts := strings.Split(image.RepoDigests[0], "@")
			imageName = truncateString(imageNameParts[0], repositoryWidth)
			tag = truncateString(getTag(image.RepoTags[0]), tagWidth)
		} else if len(image.RepoTags) > 0 {
			imageNameParts := strings.Split(image.RepoTags[0], ":")
			imageName = truncateString(imageNameParts[0], repositoryWidth)
			tag = truncateString(imageNameParts[1], tagWidth)
		} else {
			imageName = "<none>"
			tag = "<none>"
		}
		imageID := truncateString(image.ID[7:17], imageIDWidth)
		created := truncateString(time.Unix(image.Created, 0).Format("2006-01-02 15:04:05"), createdWidth)
		size := truncateString(formatSize(image.Size), sizeWidth)

		fmt.Printf("   %-*s %-*s %-*s %-*s %-*s\n",
			repositoryWidth, imageName,
			tagWidth, tag,
			imageIDWidth, imageID,
			createdWidth, created,
			sizeWidth, size,
		)
		fmt.Println(strings.Repeat("-", repositoryWidth+tagWidth+imageIDWidth+createdWidth+sizeWidth+4))
	}

	return nil
}

func getTag(repoTag string) string {
	parts := strings.SplitN(repoTag, ":", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return "latest"
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
