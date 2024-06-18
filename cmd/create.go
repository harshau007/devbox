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
	Name       string
	Package    string
	FolderPath string
	Template   string
	Port       string
}

var packageChoices = []string{"NodeLTS", "Node18", "Node20", "Python", "Rust", "Go", "Java8", "Java11", "Java17", "Java20", "Java21"}
var templateChoices = []string{"next-js", "next-ts", "nest"}

type model struct {
	cursor  int
	choice  string
	choices []string
	mode    string
}

var Dir = "Desktop"

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
			m.choice = m.choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("Which %s would you like?\n\n", m.mode))

	for i := 0; i < len(m.choices); i++ {
		if m.cursor == i {
			s.WriteString("(*) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(m.choices[i])
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
		tempFlag, _ := cmd.Flags().GetString("template")
		portFlag, _ := cmd.Flags().GetString("port")

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
				fmt.Println("Cloning repository " + name + " to " + os.Getenv("HOME") + "/" + Dir + "/" + name)

				err = os.MkdirAll(os.Getenv("HOME")+"/"+Dir+"/"+name, 0755)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				cmd := exec.Command("git", "clone", urlFlag, os.Getenv("HOME")+"/"+Dir+"/"+name)
				_, err = cmd.CombinedOutput()

				if err != nil {
					fmt.Printf("error executing the script: %v", err)
					os.Exit(1)
				}

				volFlag = Dir + "/" + name
			} else {
				fmt.Print("Enter folder path: $HOME/")
				fmt.Scanln(&volFlag)
			}
		}

		if pkgFlag == "" && tempFlag == "" {
			var selection string
			fmt.Print("Do you want to proceed with a package or a template? (package/template): ")
			fmt.Scanln(&selection)

			if strings.ToLower(selection) == "package" {
				p := tea.NewProgram(model{choices: packageChoices, mode: "package"})

				m, err := p.Run()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				if m, ok := m.(model); ok && m.choice != "" {
					pkgFlag = strings.ToLower(m.choice)
				}
			} else if strings.ToLower(selection) == "template" {
				p := tea.NewProgram(model{choices: templateChoices, mode: "template"})

				m, err := p.Run()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				if m, ok := m.(model); ok && m.choice != "" {
					tempFlag = strings.ToLower(m.choice)
				}
			} else {
				fmt.Println("Invalid selection")
				os.Exit(1)
			}
		}

		if portFlag == "" {
			fmt.Print("\nEnter the port you want to expose (default 3000): ")
			_, err := fmt.Scanln(&portFlag)
			if err != nil || portFlag == "" {
				portFlag = "3000"
			}
		}

		var container CreateContainer = CreateContainer{
			Name:       strings.ToLower(nameFlag),
			Package:    strings.ToLower(pkgFlag),
			FolderPath: volFlag,
			Template:   strings.ToLower(tempFlag),
			Port:       portFlag,
		}
		fmt.Println(container)
		_, err := createCodeInstance(container)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(createcmd)
	createcmd.PersistentFlags().StringP("name", "n", "", "Name of the container")
	createcmd.PersistentFlags().StringP("package", "p", "", "Name of the package")
	createcmd.PersistentFlags().StringP("volume", "v", "", "Path to the volume")
	createcmd.PersistentFlags().StringP("url", "u", "", "URL to the repository")
	createcmd.PersistentFlags().StringP("template", "t", "", "Template")
	createcmd.PersistentFlags().StringP("port", "P", "", "Port to expose")
}

func createCodeInstance(container CreateContainer) (string, error) {
	errCh := make(chan error)
	containerCh := make(chan string)

	go func() {
		var output []byte
		var err error
		var isExists bool
		if strings.Contains(container.Template, "next-js") {
			isExists, err = exists(container.FolderPath)
			if err != nil {
				fmt.Println("Folder: ", err)
			}
			if !isExists {
				os.MkdirAll(os.Getenv("HOME")+"/"+container.FolderPath, 0755)
			}
			cmd := exec.Command("portdevctl", container.Name, "nodelts", fmt.Sprintf("%s/%s", os.Getenv("HOME"), container.FolderPath), container.Port, container.Template)
			fmt.Println(cmd)
			output, err = cmd.CombinedOutput()
			if err != nil {
				errCh <- fmt.Errorf("error executing the script: %v", err)
				return
			}
		} else if strings.Contains(container.Template, "next-ts") {
			isExists, err = exists(container.FolderPath)
			if err != nil {
				fmt.Println("Folder: ", err)
			}
			if !isExists {
				os.MkdirAll(os.Getenv("HOME")+"/"+container.FolderPath, 0755)
			}
			cmd := exec.Command("portdevctl", container.Name, "nodelts", fmt.Sprintf("%s/%s", os.Getenv("HOME"), container.FolderPath), container.Port, container.Template)
			output, err = cmd.CombinedOutput()
			if err != nil {
				errCh <- fmt.Errorf("error executing the script: %v", err)
				return
			}
		} else if strings.Contains(container.Template, "nest") {
			isExists, err = exists(container.FolderPath)
			if err != nil {
				fmt.Println("Folder: ", err)
			}
			if !isExists {
				os.MkdirAll(os.Getenv("HOME")+"/"+container.FolderPath, 0755)
			}
			cmd := exec.Command("portdevctl", container.Name, "nodelts", fmt.Sprintf("%s/%s", os.Getenv("HOME"), container.FolderPath), container.Port, container.Template)
			output, err = cmd.CombinedOutput()
			if err != nil {
				errCh <- fmt.Errorf("error executing the script: %v", err)
				return
			}
		} else {
			isExists, err = exists(container.FolderPath)
			fmt.Println(isExists)
			if err != nil {
				fmt.Println("Folder: ", err)
			}
			if !isExists {
				os.MkdirAll(os.Getenv("HOME")+"/"+container.FolderPath, 0755)
			}
			cmd := exec.Command("portdevctl", container.Name, container.Package, fmt.Sprintf("%s/%s", os.Getenv("HOME"), container.FolderPath), container.Port, "none")
			output, err = cmd.CombinedOutput()
			if err != nil {
				errCh <- fmt.Errorf("error executing the script: %v", err)
				return
			}

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

	if !strings.HasPrefix(u.Path, "/") || !strings.Contains(u.Path, "/") {
		return "", fmt.Errorf("invalid repository URL: %s", repoURL)
	}

	pathParts := strings.Split(u.Path, "/")
	repoName := strings.TrimSuffix(pathParts[len(pathParts)-1], ".git")

	return repoName, nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
