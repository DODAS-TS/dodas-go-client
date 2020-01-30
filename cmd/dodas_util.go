package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// CreateInf is a wrapper for Infrastructure creation
func CreateInf(imURL string, templateFile string, clientConf Conf) error {

	fmt.Printf("Template: %v \n", string(templateFile))
	template, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}

	authHeader := PrepareAuthHeaders(clientConf)

	request := Request{
		URL:         string(clientConf.Im.Host),
		RequestType: "POST",
		Headers: map[string]string{
			"Authorization": authHeader,
			"Content-Type":  "text/yaml",
		},
		Content: []byte(template),
	}

	body, statusCode, err := MakeRequest(request)
	if err != nil {
		return err
	}

	if statusCode == 200 {
		stringSplit := strings.Split(string(body), "/")
		fmt.Println("InfrastructureID: ", stringSplit[len(stringSplit)-1])
	} else {
		return fmt.Errorf("Error code %d: %s", statusCode, body)
	}

	return nil
	// TODO: create .dodas dir and save infID

}

// DestroyInf is a wrapper for Infrastructure creation
func DestroyInf(imURL string, infID string, clientConf Conf) error {
	authHeader := PrepareAuthHeaders(clientConf)

	request := Request{
		URL:         imURL + "/" + infID,
		RequestType: "DELETE",
		Headers: map[string]string{
			"Authorization": authHeader,
			"Content-Type":  "text/yaml",
		},
	}

	body, statusCode, err := MakeRequest(request)
	if err != nil {
		return err
	}

	if statusCode == 200 {
		fmt.Println("Removed infrastracture ID: ", infID)
	} else {
		fmt.Println("ERROR:\n", string(body))
		return err
	}

	return nil
}
