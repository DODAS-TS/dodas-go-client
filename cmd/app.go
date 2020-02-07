package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dodas-ts/dodas-go-client/cmd/apps"
)

// MakeApps ..
func MakeApps() *cobra.Command {
	var command = &cobra.Command{
		Use:   "app",
		Short: "Install Kubernetes apps from helm charts or YAML files",
		Long: `Install Kubernetes apps from helm charts or YAML files using the "install" 
command. Helm 2 is used by default unless a --helm3 flag exists and is passed. 
You can also find the post-install message for each app with the "info" 
command.`,
		Example: `  k3sup app install
  k3sup app info inlets-operator`,
		SilenceUsage: false,
	}

	var install = &cobra.Command{
		Use:   "install",
		Short: "Install a Kubernetes app",
		Example: `  k3sup app install [APP]
  k3sup app install openfaas --help
  k3sup app install inlets-operator --token-file $HOME/do
  k3sup app install --help`,
		SilenceUsage: true,
	}

	install.PersistentFlags().String("kubeconfig", "kubeconfig", "Local path for your kubeconfig file")

	install.RunE = func(command *cobra.Command, args []string) error {

		if len(args) == 0 {
			//fmt.Printf("You can install: %s\n%s\n\n", strings.TrimRight("\n - "+strings.Join(getApps(), "\n - "), "\n - "),
			//	`Run k3sup app install NAME --help to see configuration options.`)
			return nil
		}

		return nil
	}

	command.AddCommand(install)
	install.AddCommand(apps.MakeInstallCOD(&clientConf))
	//install.AddCommand(apps.MakeInstallHTCondor(&clientConf))
	//install.AddCommand(apps.MakeInstallCMSWn(&clientConf))
	//install.AddCommand(apps.MakeInstallSparkAndJupyter(&clientConf))

	return command
}

func init() {
	cmdApps := MakeApps()
	rootCmd.AddCommand(cmdApps)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
