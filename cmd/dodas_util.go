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
	Outputs map[string]string `json:"outputs"`
}

// StatusStruct ...
type StatusStruct struct {
	Status string `json:"contmsg"`
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
func GetInfOutputs(imURL string, infID string, clientConf Conf) (outputs map[string]string, err error) {
	authHeader := PrepareAuthHeaders(clientConf)

	request := Request{
		URL:         imURL + "/" + infID + "/outputs",
		RequestType: "GET",
		Headers: map[string]string{
			"Authorization": authHeader,
			"Accept":        "application/json",
		},
	}

	body, statusCode, err := MakeRequest(request)
	if err != nil {
		return map[string]string{}, err
	}

	if statusCode != 200 {
		fmt.Println("ERROR:\n", string(body))
		return map[string]string{}, err
	}

	var bodyJSON OutputsStruct

	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		return map[string]string{}, err
	}

	return bodyJSON.Outputs, nil
}

// GetInfVMStates get ...
func GetInfVMStates(imURL string, infID string, vm string, clientConf Conf) (status string, err error) {
	authHeader := PrepareAuthHeaders(clientConf)

	request := Request{
		URL:         imURL + "/" + infID + "/vms/" + vm + "/contmsg",
		RequestType: "GET",
		Headers: map[string]string{
			"Authorization": authHeader,
			"Accept":        "application/json",
		},
	}

	body, statusCode, err := MakeRequest(request)
	if err != nil {
		return "", err
	}

	if statusCode != 200 {
		fmt.Println("ERROR:\n", string(body))
		return "", err
	}

	var bodyJSON StatusStruct

	//fmt.Println(string(body))
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		return "", err
	}

	return bodyJSON.Status, nil
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
			"Accept":        "application/json",
		},
		Content: []byte(template),
	}

	body, statusCode, err := MakeRequest(request)
	if err != nil {
		return err
	}

	if statusCode == 200 {
		fmt.Println(string(body))
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
