package apps

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/Masterminds/sprig"
	"github.com/spf13/cobra"

	"github.com/dodas-ts/dodas-go-client/pkg/utils"
)

// MakeInstallCOD ..
func MakeInstallCOD(clientConf *utils.Conf) *cobra.Command {
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
		vomsFile := Voms{
			Name: "vomsCMS",
			Content: `
asdasdad
sadasdad`,
		}

		cert := `
mycert`

		key := `
mykey`

		var params = CoDParams{
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

		infID, err := clientConf.CreateInf(b.Bytes())
		if err != nil {
			return err
		}

		fmt.Printf("InfID: %s", infID)

		return nil
	}

	return minio
}
