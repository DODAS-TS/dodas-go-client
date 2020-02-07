package cmd

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/spf13/cobra"

	"github.com/dodas-ts/dodas-go-client/apps"
)

// MakeInstallCOD ..
func MakeInstallCOD() *cobra.Command {
	var minio = &cobra.Command{
		Use:          "minio",
		Short:        "Install minio",
		Long:         `Install minio`,
		Example:      `  k3sup app install minio`,
		SilenceUsage: true,
	}

	minio.Flags().Bool("update-repo", true, "Update the helm repo")
	minio.Flags().String("access-key", "", "Provide an access key to override the pre-generated value")
	minio.Flags().String("secret-key", "", "Provide a secret key to override the pre-generated value")
	minio.Flags().Bool("distributed", false, "Deploy Minio in Distributed Mode")
	minio.Flags().String("namespace", "default", "Kubernetes namespace for the application")
	minio.Flags().Bool("persistence", false, "Enable persistence")
	minio.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set persistence.enabled=true)")

	minio.RunE = func(command *cobra.Command, args []string) error {
		// .... https://github.com/alexellis/k3sup/blob/master/cmd/apps/minio_app.go
		// https://github.com/alexellis/k3sup/blob/master/cmd/apps/kubernetes_exec.go

		// Get path
		vomsFile := apps.Voms{
			Name: "vomsCMS",
			Content: `
asdasdad
sadasdad`,
		}

		cert := `
mycert`

		key := `
mykey`

		var params = apps.CoDParams{
			SlaveNum:   2,
			ExternalIP: "{{ externalIp }}",
			VomsFile:   vomsFile,
			CacheCert:  cert,
			CacheKey:   key,
			Redirector: "myredirector.com",
		}

		// Create a new template and parse the letter into it.
		t := template.Must(template.New("CoDTemplate").Funcs(sprig.FuncMap()).Parse(CoDTemplate))

		var b bytes.Buffer
		err := t.Execute(&b, params)
		if err != nil {
			return fmt.Errorf("Failed to compile the template: %s", err)
		}

		fmt.Println(b.String())

		err = clientConf.Validate(b.Bytes())
		if err != nil {
			return err
		}

		return nil
	}

	return minio
}

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
	install.AddCommand(apps.MakeInstallCOD(clientConf))

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
