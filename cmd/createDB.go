/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// createDBCmd represents the createDB command
var createDBCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		dbuser, _ := cmd.Flags().GetString("dbuser")
		dbpass, _ := cmd.Flags().GetString("dbpass")
		dbname, _ := cmd.Flags().GetString("dbname")
		auth, _ := cmd.Flags().GetString("auth")
		pass, _ := cmd.Flags().GetString("pass")

		if dbuser == "" {
			fmt.Print("Enter the name of the Database Username: ")
			fmt.Scanln(&dbuser)
		}

		if dbpass == "" {
			fmt.Print("Enter the name of the Database Password: ")
			fmt.Scanln(&dbpass)
		}

		if dbname == "" {
			fmt.Print("Enter the name of the Databse Name: ")
			fmt.Scanln(&dbname)
		}

		if auth == "" {
			fmt.Print("Enter the name of the Mongo Express Auth: ")
			fmt.Scanln(&auth)
		}

		if pass == "" {
			fmt.Print("Enter the name of the Mongo Express Password: ")
			fmt.Scanln(&pass)
		}

		create, err := createDatabase(dbuser, dbpass, dbname, auth, pass)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(create)
	},
}

func init() {
	createDBCmd.PersistentFlags().StringP("dbuser", "", "", "DB Username")
	createDBCmd.PersistentFlags().StringP("dbpass", "", "", "DB Password")
	createDBCmd.PersistentFlags().StringP("dbname", "", "", "DB Name")
	createDBCmd.PersistentFlags().StringP("auth", "", "", "Mongo Express Username")
	createDBCmd.PersistentFlags().StringP("pass", "", "", "Mongo Express Password")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createDBCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createDBCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// DBUSER=harsh DBPASS=harsh DBNAME=test BAUTH=myadm BPASS=mypass docker compose -f /usr/local/share/devbox/MongoDB/docker-compose.yml up -d
func createDatabase(dbuser, dbpass, dbname, bauth, bpass string) (string, error) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("DBUSER=%s DBPASS=%s DBNAME=%s BAUTH=%s BPASS=%s docker compose -f /usr/local/share/devbox/MongoDB/docker-compose.yml up -d", dbuser, dbpass, dbname, bauth, bpass))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error " + err.Error())
		return "", err
	}
	if output != nil {
		return "Success", nil
	}
	return "Unsuccess", nil
}

// docker compose -f /usr/local/share/devbox/MongoDB/docker-compose.yml down -v
