/*
Copyright © 2024 Harsh Upadhyay amanupadhyay2004@gmail.com
*/
package cmd

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	newspinner "github.com/briandowns/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type CreateContainer struct {
    Name     string
    Package  string
    FolderPath string
}

var choices = []string{"NodeLTS", "Node18", "Node20", "Python", "Rust", "Go", "Java8", "Java11", "Java17", "Java20", "Java21"}

type model struct {
	cursor int
	choice string
}

var tempPath = "Desktop"

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

var createcmd = &cobra.Command{
	Use:   "create",
	Short: "Create code instance",
	Run: func(cmd *cobra.Command, args []string) {
		nameFlag, _ := cmd.Flags().GetString("name")
		pkgFlag, _ := cmd.Flags().GetString("package")
		volFlag, _ := cmd.Flags().GetString("volume")
		urlFlag, _ := cmd.Flags().GetString("url")

		if nameFlag == "" {
			fmt.Print("Enter the name of the container: ")
			fmt.Scanln(&nameFlag)
		}

		if volFlag == "" {
			if urlFlag != "" {
				name, err := getRepoName(urlFlag)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("Cloning repository " + name + " to " + os.Getenv("HOME") + tempPath + "/" + name)

				err = os.MkdirAll(os.Getenv("HOME")+ "/" + tempPath + "/" + name, 0755)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				cmd := exec.Command("git", "clone", urlFlag, os.Getenv("HOME")+ "/" + tempPath + "/" + name)  
				_, err = cmd.CombinedOutput()

				if err != nil {
					fmt.Printf("error executing the script: %v", err)
					os.Exit(1)
				}

				volFlag = tempPath + "/" + name
			} else {
				fmt.Print("Enter folder path: $HOME/")
				fmt.Scanln(&volFlag)
			}
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

func init() {
	createcmd.PersistentFlags().StringP("name", "n", "","Name of the container")
	createcmd.PersistentFlags().StringP("package", "p", "","Name of the package")
	createcmd.PersistentFlags().StringP("volume", "v", "","Path to the volume")
	createcmd.PersistentFlags().StringP("url", "u", "","URL to the repository")
}

func createCodeInstance(container CreateContainer) (string, error) {
    errCh := make(chan error)
    containerCh := make(chan string)

    go func() {
        cmd := exec.Command("portdevctl", container.Name, container.Package, fmt.Sprintf("%s/%s", os.Getenv("HOME"), container.FolderPath))
        output, err := cmd.CombinedOutput()
        if err != nil {
            errCh <- fmt.Errorf("error executing the script: %v", err)
            return
        }

        outputLines := strings.Split(strings.TrimSpace(string(output)), "\n")
        containerId := outputLines[len(outputLines)-1]

        containerCh <- containerId[:10]
    }()

	fmt.Print("\n\t")
    s := newspinner.New(newspinner.CharSets[9], 100*time.Millisecond)
    s.Writer = os.Stderr
    s.Suffix = " Creating Code Instance...\n"
    s.Start()
    defer s.Stop()

    select {
    case err := <-errCh:
        s.Stop()
        return "", err
    case containerId := <-containerCh:
        s.Stop()
        return containerId, nil
    }
}

func getRepoName(repoURL string) (string, error) {
    u, err := url.Parse(repoURL)
    if err != nil {
        return "", err
    }

    host := u.Hostname()
    if !strings.HasSuffix(host, "github.com") && !strings.HasSuffix(host, "gitlab.com") {
        return "", fmt.Errorf("invalid repository URL: %s", repoURL)
    }

    if !strings.HasPrefix(u.Path, "/") {
        return "", fmt.Errorf("invalid repository URL: %s", repoURL)
    }

    pathComponents := strings.Split(u.Path, "/")

    if len(pathComponents) < 3 {
        return "", fmt.Errorf("invalid repository URL: %s", repoURL)
    }

    if strings.HasSuffix(host, "github.com") {
        return pathComponents[2], nil
    } else {
        return pathComponents[3], nil
    }
}