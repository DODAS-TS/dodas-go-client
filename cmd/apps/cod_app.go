package apps

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/spf13/cobra"

	"github.com/dodas-ts/dodas-go-client/pkg/utils"
)

// MakeInstallCOD ..
func MakeInstallCOD(clientConf *utils.Conf) *cobra.Command {
	var cod *cobra.Command = &cobra.Command{
		Use:          "cod",
		Short:        "Install CachingOnDemand",
		Long:         ``,
		Example:      `dodas app install cod`,
		SilenceUsage: true,
	}

	cod.Flags().Int("nslaves", 2, "Number of cache servers")
	cod.Flags().Bool("gsi-enabled", false, "Enable GSI-VO auth")
	cod.Flags().String("origin-host", "xrootd-cms.infn.it", "Origin server host")
	cod.Flags().Int("origin-port", 1094, "Origin server port")
	cod.Flags().String("redirector-host", "", "Cache redirector host")
	cod.Flags().Int("redirector-port", 31213, "Cache redirector port")
	cod.Flags().String("voms-file", "", "Provide a voms file for a vo (example --voms-file cms.txt=/tmp/voms/cms.txt)")
	cod.Flags().String("x509-cert", "", "path to the certificate for cache server auth with remote")
	cod.Flags().String("x509-key", "", "path to the certificate key for cache server auth with remote")
	cod.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set persistence.enabled=true)")

	cod.RunE = func(command *cobra.Command, args []string) error {

		gsiEnabled, _ := cod.Flags().GetBool("gsi-enabled")
		vomsFileConf, _ := cod.Flags().GetString("voms-file")
		certFileConf, _ := cod.Flags().GetString("x509-cert")
		keyFileConf, _ := cod.Flags().GetString("x509-key")
		redirectorHost, _ := cod.Flags().GetString("redirector-host")
		originHost, _ := cod.Flags().GetString("origin-host")
		redirectorPort, _ := cod.Flags().GetInt("redirector-port")
		originPort, _ := cod.Flags().GetInt("origin-port")
		nslaves, _ := cod.Flags().GetInt("nslaves")

		var params CoDParams
		var slaves = SlavesStruct{
			SlaveNum: nslaves,
		}

		if gsiEnabled {
			vomsFileStrings := strings.Split(vomsFileConf, "=")
			if len(vomsFileStrings) < 2 {
				return fmt.Errorf("Invalid voms file format. Please use e.g.: --voms-file cms.txt=/tmp/voms/cms.txt")
			}

			vomsContent, err := ioutil.ReadFile(vomsFileStrings[1])
			if err != nil {
				return fmt.Errorf("Failed to read voms file: %s", err)
			}

			vomsFile := Voms{
				Name:    vomsFileStrings[0],
				Content: string(vomsContent),
			}

			cert, err := ioutil.ReadFile(certFileConf)
			if err != nil {
				return fmt.Errorf("Failed to read X509 cert file: %s", err)
			}
			key, err := ioutil.ReadFile(keyFileConf)
			if err != nil {
				return fmt.Errorf("Failed to read X509 cert key file: %s", err)
			}

			params = CoDParams{
				Slaves:         slaves,
				ExternalIP:     "{{ externalIp }}",
				VomsFile:       vomsFile,
				GsiEnabled:     gsiEnabled,
				CacheCert:      string(cert),
				CacheKey:       string(key),
				RedirectorHost: redirectorHost,
				RedirectorPort: redirectorPort,
				OriginHost:     originHost,
				OriginPort:     originPort,
			}
		}

		params = CoDParams{
			Slaves:         slaves,
			GsiEnabled:     gsiEnabled,
			ExternalIP:     "{{ externalIp }}",
			RedirectorHost: redirectorHost,
			RedirectorPort: redirectorPort,
			OriginHost:     originHost,
			OriginPort:     originPort,
		}

		// Create a new template and parse the letter into it.
		t := template.Must(template.New("CoDTemplate").Funcs(sprig.FuncMap()).Parse(CoDTemplate))

		var b bytes.Buffer
		err := t.Execute(&b, params)
		if err != nil {
			return fmt.Errorf("Failed to compile the template: %s", err)
		}

		//fmt.Println(b.String())

		err = clientConf.Validate(b.Bytes())
		if err != nil {
			return err
		}

		infID, err := clientConf.CreateInf(b.Bytes())
		if err != nil {
			return err
		}

		tempFile, err := ioutil.TempFile("/tmp/", "cod_"+infID+"_")
		defer tempFile.Close()
		if err != nil {
			return err
		}

		fmt.Printf("Writing compiled template in: %s", tempFile.Name())
		_, err = tempFile.Write(b.Bytes())
		if err != nil {
			return err
		}

		return nil
	}

	return cod
}
