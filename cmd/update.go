package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update <infID> <template>",
	Short: "Update the number of vms to satisfy the new template",
	Long:  ``,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update called")
		err := clientConf.Validate(args[1])
		if err != nil {
			panic(err)
		}

		fmt.Printf("Updating infID %s with: %s \n", args[0], args[1])
		err = clientConf.UpdateInf(string(clientConf.Im.Host), args[0], args[1])
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
