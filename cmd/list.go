/*
Copyright © 2024 Harsh Upadhyay amanupadhyay2004@gmail.com
*/
package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// docker ps -a --filter "label=createdBy=DevControl"

type ContainerInfo struct {
	ID     string
	Name   string
	Image  string
	Status string
	Volume string
	Url    string
}

var allFlag bool

var listCmd = &cobra.Command{
	Use:   "ps",
	Short: "List containers",
	Run: func(cmd *cobra.Command, args []string) {
		if allFlag {
			listAllContainers()
		} else {
			listContainers()
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().BoolVarP(&allFlag, "all", "a", false, "List all containers")
}

func listContainers() error {
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return err
    }
    defer cli.Close()

    filters := filters.NewArgs(filters.Arg("label", "createdBy=DevControl"))
    containers, err := cli.ContainerList(ctx, containertypes.ListOptions{Filters: filters})
    if err != nil {
        return err
    }

    const (
        containerIDWidth = 12
        imageWidth       = 20
        createdWidth     = 20
        statusWidth      = 20
        portsWidth       = 20
        namesWidth       = 20
    )

	volumeWidth := maxVolumeWidth(containers)

	if len(containers) == 0 {
        fmt.Println("No containers found.")
        return nil
    }
    fmt.Println(strings.Repeat("-", containerIDWidth+imageWidth+volumeWidth+createdWidth+statusWidth+portsWidth+namesWidth+11))

    fmt.Printf("   %-*s\t %-*s %-*s\t %-*s\t %-*s\t %-*s\t %-*s\n",
        containerIDWidth, "CONTAINER ID",
        imageWidth, "IMAGE",
        volumeWidth, "VOLUME",
        createdWidth, "CREATED",
        statusWidth, "STATUS",
        portsWidth, "URL",
        namesWidth, "NAMES")
    fmt.Println(strings.Repeat("-", containerIDWidth+imageWidth+volumeWidth+createdWidth+statusWidth+portsWidth+namesWidth+11))

    for _, container := range containers {
		var url string
		if len(container.Ports) > 0 {
			url = fmt.Sprintf("http://%s:%d", container.Ports[0].IP, container.Ports[0].PublicPort)
		}
		if url == "" {
			url = "Not Available"
		}
        containerID := truncateString(container.ID, containerIDWidth)
        image := truncateString(container.Image, imageWidth)
        volume := truncateString(container.Mounts[0].Source, volumeWidth)
        created := truncateString(time.Unix(container.Created, 0).Format("2006-01-02 15:04:05"), createdWidth)
        status := truncateString(container.Status, statusWidth)
        urls := url
        names := truncateString(container.Names[0][1:], namesWidth)

        fmt.Printf("   %-*s\t %-*s %-*s\t %-*s\t %-*s\t %-*s\t %-*s\n",
            containerIDWidth, containerID,
            imageWidth, image,
            volumeWidth, volume,
            createdWidth, created,
            statusWidth, status,
            portsWidth, urls,
            namesWidth, names,
        )
        fmt.Println(strings.Repeat("-", containerIDWidth+imageWidth+volumeWidth+createdWidth+statusWidth+portsWidth+namesWidth+11))
    }

    return nil
}

func listAllContainers() error {
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return err
    }
    defer cli.Close()

	filters := filters.NewArgs(
		filters.Arg("label", "createdBy=DevControl"),
	)
    containers, err := cli.ContainerList(ctx, containertypes.ListOptions{All: true, Filters: filters})
    if err != nil {
        return err
    }

    const (
        containerIDWidth = 12
        imageWidth       = 20
        createdWidth     = 20
        statusWidth      = 20
        portsWidth       = 20
        namesWidth       = 20
    )

	volumeWidth := maxVolumeWidth(containers)

	if len(containers) == 0 {
        fmt.Println("No containers found.")
        return nil
    }

    fmt.Println(strings.Repeat("-", containerIDWidth+imageWidth+volumeWidth+createdWidth+statusWidth+portsWidth+namesWidth+11))

    fmt.Printf("   %-*s\t %-*s %-*s\t %-*s\t %-*s\t %-*s\t %-*s\n",
        containerIDWidth, "CONTAINER ID",
        imageWidth, "IMAGE",
        volumeWidth, "VOLUME",
        createdWidth, "CREATED",
        statusWidth, "STATUS",
        portsWidth, "URL",
        namesWidth, "NAMES")
    fmt.Println(strings.Repeat("-", containerIDWidth+imageWidth+volumeWidth+createdWidth+statusWidth+portsWidth+namesWidth+11))

    for _, container := range containers {
		var url string
		if len(container.Ports) > 0 {
			url = fmt.Sprintf("http://%s:%d", container.Ports[0].IP, container.Ports[0].PublicPort)
		}
		if url == "" {
			url = "Not Available"
		}
        containerID := truncateString(container.ID, containerIDWidth)
        image := truncateString(container.Image, imageWidth)
        volume := truncateString(container.Mounts[0].Source, volumeWidth)
        created := truncateString(time.Unix(container.Created, 0).Format("2006-01-02 15:04:05"), createdWidth)
        status := truncateString(container.Status, statusWidth)
        urls := url
        names := truncateString(container.Names[0][1:], namesWidth)

        fmt.Printf("   %-*s\t %-*s %-*s\t %-*s\t %-*s\t %-*s\t %-*s\n",
            containerIDWidth, containerID,
            imageWidth, image,
            volumeWidth, volume,
            createdWidth, created,
            statusWidth, status,
            portsWidth, urls,
            namesWidth, names,
        )
    }

    fmt.Println(strings.Repeat("-", containerIDWidth+imageWidth+volumeWidth+createdWidth+statusWidth+portsWidth+namesWidth+11))

    return nil
}

func maxVolumeWidth(containers []types.Container) int {
    maxWidth := 0
    for _, container := range containers {
        if len(container.Mounts) > 0 {
            volumePath := container.Mounts[0].Source
            if len(volumePath) > maxWidth {
                maxWidth = len(volumePath)
            }
        }
    }
    return maxWidth
}

func truncateString(str string, width int) string {
    if len(str) > width {
        return str[:width]
    }
    return str
}
