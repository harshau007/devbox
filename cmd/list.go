/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// docker ps -a --filter "label=createdBy=instacode"

type ContainerInfo struct {
	ID     string
	Name   string
	Image  string
	Status string
	Volume string
	Url    string
}

type modelTable struct {
	table         *table.Table
	containerList []ContainerInfo
}

func (m *modelTable) Init() tea.Cmd {
	return nil
}

func (m *modelTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table = m.table.Width(msg.Width)
		m.table = m.table.Height(msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
		}
	}
	return m, cmd
}

func (m modelTable) View() string {
	return "\n" + m.table.String() + "\n(press q to quit)\n"
}


func (m *modelTable) initTable() {
	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Copy().Foreground(lipgloss.Color("252")).Bold(true)
	selectedStyle := baseStyle.Copy().Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))
	typeColors := map[string]lipgloss.Color{
		"Bug":      lipgloss.Color("#D7FF87"),
		"Electric": lipgloss.Color("#FDFF90"),
		"Fire":     lipgloss.Color("#FF7698"),
		"Flying":   lipgloss.Color("#FF87D7"),
		"Grass":    lipgloss.Color("#75FBAB"),
		"Ground":   lipgloss.Color("#FF875F"),
		"Normal":   lipgloss.Color("#929292"),
		"Poison":   lipgloss.Color("#7D5AFC"),
		"Water":    lipgloss.Color("#00E2C7"),
	}
	dimTypeColors := map[string]lipgloss.Color{
		"Bug":      lipgloss.Color("#97AD64"),
		"Electric": lipgloss.Color("#FCFF5F"),
		"Fire":     lipgloss.Color("#BA5F75"),
		"Flying":   lipgloss.Color("#C97AB2"),
		"Grass":    lipgloss.Color("#59B980"),
		"Ground":   lipgloss.Color("#C77252"),
		"Normal":   lipgloss.Color("#727272"),
		"Poison":   lipgloss.Color("#634BD0"),
		"Water":    lipgloss.Color("#439F8E"),
	}
	columns := []string{
		"ID",
		"Name",
		"Image",
		"Status",
		"Volume",
		"URL",
	}

	rows := make([][]string, len(m.containerList))
	for i, container := range m.containerList {
		rows[i] = []string{container.ID, container.Name, container.Image, container.Status, container.Volume, container.Url}
	}

	m.table = table.New().Headers(columns...).Rows(rows...).Height(len(m.containerList)+2).Border(lipgloss.NormalBorder()).
	BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
	StyleFunc(func(row, col int) lipgloss.Style {
		if row == 0 {
			return headerStyle
		}

		if rows[row-1][1] == "Pikachu" {
			return selectedStyle
		}

		even := row%2 == 0

		switch col {
		case 2, 3: 
			c := typeColors
			if even {
				c = dimTypeColors
			}

			color := c[fmt.Sprint(rows[row-1][col])]
			return baseStyle.Copy().Foreground(color)
		}

		if even {
			return baseStyle.Copy().Foreground(lipgloss.Color("252"))
		}
		return baseStyle.Copy().Foreground(lipgloss.Color("252"))
	}).
	Border(lipgloss.ThickBorder())

}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cont, err := listContainers()
		if err != nil {
			fmt.Println(err)
			return
		}

		m := modelTable{containerList: cont}
		m.initTable()
		p := tea.NewProgram(&m, tea.WithAltScreen())
		_, err = p.Run()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listContainers() ([]ContainerInfo, error) {
	containerlist := []ContainerInfo{}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	filters := filters.NewArgs(filters.Arg("label", "createdBy=instacode"))
	containers, err := cli.ContainerList(ctx, containertypes.ListOptions{All: true, Filters: filters})
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		var url string
		if len(container.Ports) > 0 {
			url = fmt.Sprintf("http://%s:%d", container.Ports[0].IP, container.Ports[0].PublicPort)
		}
		containerlist = append(containerlist, ContainerInfo{
			ID:     container.ID[:10],
			Name:   container.Names[0][1:],
			Image:  container.Image,
			Status: container.Status,
			Url:    url,
			Volume: container.Mounts[0].Source,
		})
	}

	return containerlist, nil
}
