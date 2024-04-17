/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

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
			fmt.Print("\nEnter the name of the container: ")
			fmt.Scanln(&nameFlag)
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
	stopCmd.PersistentFlags().StringP("name", "n", "","Name of the container")
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
	return "Stopped container: " + contInfo.ID, nil
}
