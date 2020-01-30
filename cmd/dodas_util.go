package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/dciangot/toscalib"
)

// OutputsStruct ...
type OutputsStruct struct {
	Outputs []string `json:"outputs"`
}

// CreateInf is a wrapper for Infrastructure creation
func CreateInf(imURL string, templateFile string, clientConf Conf) (infID string, err error) {

	fmt.Printf("Template: %v \n", string(templateFile))
	template, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return "", err
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
		return "", err
	}

	if statusCode == 200 {
		stringSplit := strings.Split(string(body), "/")
		fmt.Println("InfrastructureID: ", stringSplit[len(stringSplit)-1])
	} else {
		return "", fmt.Errorf("Error code %d: %s", statusCode, body)
	}

	stringSplit := strings.Split(string(body), "/")
	return stringSplit[len(stringSplit)-1], nil
	// TODO: create .dodas dir and save infID

}

// GetInfOutputs get ...
func GetInfOutputs(imURL string, infID string, clientConf Conf) (outputs []string, err error) {
	authHeader := PrepareAuthHeaders(clientConf)

	request := Request{
		URL:         imURL + "/" + infID + "/outputs",
		RequestType: "Get",
		Headers: map[string]string{
			"Authorization": authHeader,
			"Content-Type":  "application/json",
		},
	}

	body, statusCode, err := MakeRequest(request)
	if err != nil {
		return []string{}, err
	}

	if statusCode == 200 {
		fmt.Println("Outputs: ", string(body))
	} else {
		fmt.Println("ERROR:\n", string(body))
		return []string{}, err
	}

	var bodyJSON OutputsStruct

	json.Unmarshal([]byte(body), &bodyJSON)

	return bodyJSON.Outputs, nil
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

// UpdateInf ..
func UpdateInf(imURL string, infID string, templateFile string, clientConf Conf) error {

	fmt.Printf("Template: %v \n", string(templateFile))
	template, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}

	authHeader := PrepareAuthHeaders(clientConf)

	request := Request{
		URL:         string(clientConf.Im.Host) + "/" + infID,
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

// Validate TOSCA template
func Validate(templateFile string) error {
	fmt.Println("validate called")
	var t toscalib.ServiceTemplateDefinition
	template, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}

	err = t.Parse(bytes.NewBuffer(template))
	if err != nil {
		fmt.Printf("ERROR: Invalid template for %v", err)
		return err
	}
	// t.TopologyTemplate.NodeTemplates

	//t.TopologyTemplate.NodeTemplates["Type"]

	//typeList := make(map[string][]string)

	inputs := make(map[string][]string)
	templs := make(map[string][]string)

	for name := range t.TopologyTemplate.NodeTemplates {
		//fmt.Println(name)

		for templ := range t.TopologyTemplate.NodeTemplates[name].Properties {
			//fmt.Println(templ)
			value := t.TopologyTemplate.NodeTemplates[name].Properties[templ].Value
			ft := t.TopologyTemplate.NodeTemplates[name].Properties[templ].Function
			if value != "" && value != nil || ft != "" {
				templs[name] = append(templs[name], templ)
			}
			//fmt.Print("-----\n")
		}

		//fmt.Print("++++\n")
		derived := t.NodeTypes[t.TopologyTemplate.NodeTemplates[name].Type].DerivedFrom
		for derived != "" {
			for interf := range t.NodeTypes[derived].Properties {
				//fmt.Println(interf)
				inputs[name] = append(inputs[name], interf)
			}
			//fmt.Println(derived)
			derived = t.NodeTypes[derived].DerivedFrom
		}

		for interf := range t.NodeTypes[t.TopologyTemplate.NodeTemplates[name].Type].Properties {
			inputs[name] = append(inputs[name], interf)
		}

	}
	//fmt.Println(inputs)
	//fmt.Println(templs)

	for node := range templs {
		//fmt.Println(node)
		for nodeParam := range templs[node] {
			isPresent := false
			for param := range inputs[node] {

				if inputs[node][param] == templs[node][nodeParam] {
					isPresent = true
				}
			}
			//fmt.Printf("%v %v\n", templs[node][nodeParam], isPresent)
			if !isPresent {
				fmt.Printf("%v not defined in type %v \n", templs[node][nodeParam], t.TopologyTemplate.NodeTemplates[node].Type)
				return fmt.Errorf("ERROR: Invalid template for %v", node)
			}
		}
		//fmt.Print("-----\n")
	}

	fmt.Print("Template OK\n")
	return nil
}
