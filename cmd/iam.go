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
	Short: "Initialize a IAM context for automatic token refresh",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		isTokenUsed := (clientConf.Im.Token != "" || clientConf.Cloud.AuthVersion == "3.x_oidc_access_token")
		isRefreshSet := clientConf.AllowRefresh.IAMTokenEndpoint != ""

		if !isTokenUsed && !isRefreshSet {
			panic(fmt.Errorf("Token not used anywhere in config or refresh endpoint missing"))
		}

		fmt.Println("Removing old dump files and creating new ones")

		// Remove dumps files if exists
		if _, err := os.Stat(clientConf.AllowRefresh.AccessTokenFile); err == nil {
			err := os.Remove(clientConf.AllowRefresh.AccessTokenFile)
			if err != nil {
				panic(err)
			}
		}

		if _, err := os.Stat(clientConf.AllowRefresh.RefreshTokenFile); err == nil {
			err = os.Remove(clientConf.AllowRefresh.RefreshTokenFile)
			if err != nil {
				panic(err)
			}
		}

		if clientConf.Im.Token == "" {
			panic(fmt.Errorf("Error: access token not specified to IM"))
		}

		if err := ioutil.WriteFile(clientConf.AllowRefresh.AccessTokenFile, []byte(clientConf.Im.Token), os.FileMode(int(0600))); err != nil {
			log.Fatal(err)
		}

		token, err := clientConf.GetRefreshToken()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Got refresh token: %s", token)

		if err := ioutil.WriteFile(clientConf.AllowRefresh.RefreshTokenFile, []byte(token), os.FileMode(int(0600))); err != nil {
			log.Fatal(err)
		}

	},
}

// iamCmd represents the iam command
var iamCmd = &cobra.Command{
	Use:   "iam",
	Short: "Wrapper command for IAM interaction",
	Long:  ``,
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
