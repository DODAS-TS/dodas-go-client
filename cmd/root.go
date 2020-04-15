package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"

	"github.com/dodas-ts/dodas-go-client/pkg/utils"
)

var (
	version      bool
	cfgFile      string
	templateFile string
	infID        string
	//clientset    *kubernetes.Clientset
	kubeconfig string
	clientConf utils.Conf
)

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//fmt.Println(home)

		// Search config in home directory with name ".dodas_go_client" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dodas")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
		clientConf = clientConf.GetConf(viper.ConfigFileUsed())

		//fmt.Printf("TOKEN: %s", clientConf.Im.Token)
		//if clientConf.im.Password == "" {
		//	fmt.Println("No password")
		//}
	}

}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dodas",
	Short: "A self-sufficient client for DODAS deployments",
	Long: `A self-sufficient client for DODAS deployments.
Default configuration file searched in $HOME/.dodas.yaml

Usage examples:
"""
# CREATE A CLUSTER FROM TEMPLATE
dodas create --template my_tosca_template.yml

# VALIDATE TOSCA TEMPLATE
dodas validate --template my_tosca_template.yml
"""`,

	Run: func(cmd *cobra.Command, args []string) {
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// VersionString ..
var VersionString string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	VersionString = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// BuildDoc ...
func BuildDoc() {
	err := doc.GenMarkdownTree(rootCmd, "docs")
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "DODAS config file (default is $HOME/.dodas.yaml)")
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "/etc/kubernetes/admin.conf", "Kubernetes config file (default is /etc/kubernetes/admin.conf)")
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "DODAS client version")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
