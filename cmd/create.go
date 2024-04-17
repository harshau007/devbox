/*
Copyright © 2024 Harsh Upadhyay amanupadhyay2004@gmail.com
*/
package cmd

import (	
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	tea "github.com/charmbracelet/bubbletea"
)

type CreateContainer struct {
    Name     string
    Package  string
    FolderPath string
}

var choices = []string{"NodeLTS", "Node18", "Node20", "Python", "Rust"}

type model struct {
	cursor int
	choice string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.choice = choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString("\nWhich package would you like?\n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(*) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

var createcmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		nameFlag, _ := cmd.Flags().GetString("name")
		pkgFlag, _ := cmd.Flags().GetString("package")
		volFlag, _ := cmd.Flags().GetString("volume")

		if nameFlag == "" {
			fmt.Print("\nEnter the name of the container: ")
			fmt.Scanln(&nameFlag)
		}

		if volFlag == "" {
			fmt.Print("\nEnter folder path: $HOME/")
			fmt.Scanln(&volFlag)
		}

		if pkgFlag == "" {
			p := tea.NewProgram(model{})

			m, err := p.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		
			if m, ok := m.(model); ok && m.choice != "" {
				pkgFlag = strings.ToLower(m.choice)
			}
		}

		var container CreateContainer = CreateContainer{
			Name:     strings.ToLower(nameFlag),
			Package:  strings.ToLower(pkgFlag),
			FolderPath: volFlag,
		}
		_, err := createCodeInstance(container)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func createCodeInstance(container CreateContainer) (string, error) {
    cmd := exec.Command("./port", container.Name, container.Package, fmt.Sprintf("%s/%s", os.Getenv("HOME"), container.FolderPath))
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Error executing the script: %v\n", err)
        return "", err
    }

    outputLines := strings.Split(strings.TrimSpace(string(output)), "\n")
    containerId := outputLines[len(outputLines)-1]
    fmt.Printf("\nContainer created with ID: %s\n", containerId)

    return containerId[:10], nil
}