/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

type containerDetail struct {
	Name string
}

type modelContainer struct {
	cursor int
	choice string
}

var choicesContainers = []string{}

func (m modelContainer) Init() tea.Cmd {
	return nil
}

func (m modelContainer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.choice = choicesContainers[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choicesContainers) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choicesContainers) - 1
			}
		}
	}

	return m, nil
}

func (m modelContainer) View() string {
	s := strings.Builder{}
	s.WriteString("\nWhich container would you like to stop?\n\n")

	for i := 0; i < len(choicesContainers); i++ {
		if m.cursor == i {
			s.WriteString("(*) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choicesContainers[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		nameFlag, _ := cmd.Flags().GetString("name")
		if nameFlag == "" {
			containerInfo, err := containerNames()
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, container := range containerInfo {
				choicesContainers = append(choicesContainers, container.Name)
			}

			if len(choicesContainers) == 0 {
				fmt.Println("No running containers found!")
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
		msg, err := stopcontainer(nameFlag)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(msg)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.PersistentFlags().StringP("name", "n", "", "Name of the container")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func stopcontainer(name string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer cli.Close()

	contInfo, err := cli.ContainerInspect(ctx, name)
	if err != nil {
		return "", err
	}
	if contInfo.State.Running {
		err = cli.ContainerStop(ctx, name, containertypes.StopOptions{})
		if err != nil {
			return "", err
		}
	}
	return "\nStopped container: " + contInfo.ID[:10] + "\n", nil
}

func containerNames() ([]containerDetail, error) {
	containerlist := []containerDetail{}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	filters := filters.NewArgs(filters.Arg("label", "createdBy=DevControl"))
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
