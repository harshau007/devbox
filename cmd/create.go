/*
Copyright © 2024 Harsh Upadhyay amanupadhyay2004@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
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
	s.WriteString("Which package would you like?\n\n")

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

type loaderModel struct {
    spinner  spinner.Model
    quitting bool
    err      error
    id       string
    done     bool
    mutex    sync.Mutex
}

func initialModel() loaderModel {
    s := spinner.New()
    s.Spinner = spinner.Dot
    s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
    return loaderModel{spinner: s}
}

func (m *loaderModel) Init() tea.Cmd {
    return m.spinner.Tick
}

func (m *loaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            m.quitting = true
            return m, tea.Quit
        default:
            return m, nil
        }

    case spinner.TickMsg:
        var cmd tea.Cmd
        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd

    default:
        return m, nil
    }
}

func (m *loaderModel) View() string {
    if m.err != nil {
        return m.err.Error()
    }
    str := fmt.Sprintf("\n  %s Creating Code Instance... \n\n", m.spinner.View())
    if m.id != "" {
        str += fmt.Sprintf("Container created with ID: %s\n", m.id)
    }
    return str
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
			fmt.Print("Enter the name of the container: ")
			fmt.Scanln(&nameFlag)
		}

		if volFlag == "" {
			fmt.Print("Enter folder path: $HOME/")
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
    m := initialModel()
    p := tea.NewProgram(&m)
    var wg sync.WaitGroup
    wg.Add(1)
    done := make(chan struct{})

    go func() {
        defer wg.Done()
        cmd := exec.Command("port", container.Name, container.Package, fmt.Sprintf("%s/%s", os.Getenv("HOME"), container.FolderPath))
        output, err := cmd.CombinedOutput()
        if err != nil {
            m.mutex.Lock()
            m.err = fmt.Errorf("error executing the script: %v", err)
            m.mutex.Unlock()
            close(done)
            return
        }

        outputLines := strings.Split(strings.TrimSpace(string(output)), "\n")
        containerId := outputLines[len(outputLines)-1]

        m.mutex.Lock()
        m.id = containerId[:10]
        m.done = true
        m.mutex.Unlock()
        close(done)
    }()

    go func() {
        for {
            select {
            case <-done:
                p.Kill()
                return
            default:
                _, err := p.Run()
                if err != nil {
                    return
                }
            }
        }
    }()

    wg.Wait()
    return m.id, m.err
}