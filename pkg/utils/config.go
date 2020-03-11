package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

// ConfCloud ...
type ConfCloud struct {
	ID            string `yaml:"id"`
	Type          string `yaml:"type"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	Host          string `yaml:"host"`
	Tenant        string `yaml:"tenant"`
	AuthURL       string `yaml:"auth_url,omitempty"`
	AuthVersion   string `yaml:"auth_version"`
	Domain        string `yaml:"domain,omitempty"`
	ServiceRegion string `yaml:"service_region,omitempty"`
}

// ConfIM ..
type ConfIM struct {
	ID       string `yaml:"id"`
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Token    string `yaml:"token,omitempty"`
}

// TokenRefreshConf ..
type TokenRefreshConf struct {
	ClientID         string `yaml:"client_id"`
	ClientSecret     string `yaml:"client_secret"`
	IAMTokenEndpoint string `yaml:"iam_endpoint"`
	RefreshTokenFile string `yaml:"refresh_file"`
	AccessTokenFile  string `yaml:"access_file"`
}

// Conf ..
type Conf struct {
	Im           ConfIM           `yaml:"im"`
	Cloud        ConfCloud        `yaml:"cloud"`
	AllowRefresh TokenRefreshConf `yaml:"allowrefresh,omitempty"`
}

// GetConf ..
func (clientConf Conf) GetConf(path string) Conf {

	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = yaml.UnmarshalStrict(f, &clientConf)
	if err != nil {
		panic(err)
	}

	// if access token is dumped use it
	isTokenUsed := (clientConf.Im.Token != "" || clientConf.Cloud.AuthVersion == "3.x_oidc_access_token")
	isRefreshSet := clientConf.AllowRefresh.IAMTokenEndpoint != ""

	if isTokenUsed && isRefreshSet {
		tokenBytes, err := ioutil.ReadFile(clientConf.AllowRefresh.AccessTokenFile)
		if err != nil {
			fmt.Printf("Failed to read access token file %s, not going to use cache tokens: %s", clientConf.AllowRefresh.AccessTokenFile, err.Error())
			if clientConf.Im.Token != "" {
				fmt.Printf("Saving access token from configfile %s \n", clientConf.AllowRefresh.AccessTokenFile)
				if err := ioutil.WriteFile(clientConf.AllowRefresh.AccessTokenFile, []byte(clientConf.Im.Token), os.FileMode(int(0600))); err != nil {
					panic(err)
				}
			}
			return clientConf
		}

		if clientConf.Cloud.AuthVersion == "3.x_oidc_access_token" {
			clientConf.Cloud.Password = string(tokenBytes)
		}
		if clientConf.Im.Token != "" {
			clientConf.Im.Token = string(tokenBytes)
		}

		_, err = clientConf.ListInfIDs()
		if err != nil {

			re := regexp.MustCompile(`^.*OIDC auth Token expired.*`)
			if re.Match([]byte(err.Error())) {

				fmt.Printf("Token expired, trying to refresh the token ")

				clientConf, err = clientConf.GetNewToken()
				if err != nil {
					panic(err)
				}

				// Dump the new token
				fmt.Printf("Saving new access token in %s \n", clientConf.AllowRefresh.AccessTokenFile)
				if err := ioutil.WriteFile(clientConf.AllowRefresh.AccessTokenFile, []byte(clientConf.Im.Token), os.FileMode(int(0600))); err != nil {
					panic(err)
				}
			}
		}
	}
	//fmt.Printf("--- c.im:\n%v\n\n", string(c.Im.Host))

	return clientConf
}
