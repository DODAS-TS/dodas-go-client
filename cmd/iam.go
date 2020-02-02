package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")

		if clientConf.Im.Token == "" {
			panic(fmt.Errorf("Error: access token not specified to IM"))
		}

		token, err := clientConf.GetRefreshToken()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Got refresh token: %s", token)

		if err := ioutil.WriteFile(clientConf.AllowRefresh.DumpFile, []byte(token), os.FileMode(600)); err != nil {
			log.Fatal(err)
		}

		tokenBytes, err := ioutil.ReadFile(clientConf.AllowRefresh.DumpFile)
		if err != nil {
			panic(err)
		}

		accessToken, err := clientConf.GetAccessToken(string(tokenBytes))
		if err != nil {
			panic(err)
		}

		fmt.Printf("Access token: %s", accessToken)
	},
}

// iamCmd represents the iam command
var iamCmd = &cobra.Command{
	Use:   "iam",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("iam called")
	},
}

func init() {
	rootCmd.AddCommand(iamCmd)
	iamCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// iamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// iamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
