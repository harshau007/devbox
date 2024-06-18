/*
Copyright © 2024 Harsh Upadhyay amanupadhyay2004@gmail.com
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start containers",
	Run: func(cmd *cobra.Command, args []string) {
		nameFlag, _ := cmd.Flags().GetString("name")
		if nameFlag == "" {
			containerInfo, err := containerList()
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, container := range containerInfo {
				choicesContainers = append(choicesContainers, container.Name)
			}

			if len(choicesContainers) == 0 {
				fmt.Println("No stopped containers found!")
				return
			}
			p := tea.NewProgram(modelContainer{})

			m, err := p.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if m, ok := m.(modelContainer); ok && m.choice != "" {
				nameFlag = strings.ToLower(m.choice)
			}

			indexToRemove := -1
			for i, v := range choicesContainers {
				if v == nameFlag {
					indexToRemove = i
					break
				}
			}

			if indexToRemove != -1 {
				_ = append(choicesContainers[:indexToRemove], choicesContainers[indexToRemove+1:]...)
			}
		}
		_, err := startContainer(nameFlag)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringP("name", "n", "", "Name of the container")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func startContainer(name string) (string, error) {
	cmd := exec.Command("startdevctl", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing the script: %v\n", err)
		return "", err
	}

	outputLines := strings.Split(strings.TrimSpace(string(output)), "\n")
	containerId := outputLines[len(outputLines)-1]
	fmt.Printf("\nContainer created with ID: %s\n", containerId[:10])

	return "", nil
}

func containerList() ([]containerDetail, error) {
	containerlist := []containerDetail{}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	filters := filters.NewArgs(
		filters.Arg("label", "createdBy=DevBox"),
		filters.Arg("status", "exited"),
	)
	containers, err := cli.ContainerList(ctx, containertypes.ListOptions{Filters: filters})
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		containerlist = append(containerlist, containerDetail{
			Name: container.Names[0][1:],
		})
	}

	return containerlist, nil
}
