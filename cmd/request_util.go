package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// Request input struct
type Request struct {
	URL         string
	RequestType string
	Headers     map[string]string
	AuthUser    string
	AuthPwd     string
	Content     []byte
	Timeout     time.Duration
}

func validateRequest(r Request) (Request, error) {

	validatedRequest := r

	if &r.Timeout == nil {
		validatedRequest.Timeout = 5 * time.Minute
	}

	if r.URL == "" {
		return Request{}, fmt.Errorf("URL not specified")
	}

	if r.RequestType == "" {
		validatedRequest.RequestType = "GET"
	}

	return validatedRequest, nil
}

// PrepareAuthHeaders ..
func PrepareAuthHeaders(clientConf Conf) string {

	var authHeaderCloudList []string

	fields := reflect.TypeOf(clientConf.Cloud)
	values := reflect.ValueOf(clientConf.Cloud)

	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		value := values.Field(i)

		if value.Interface() != "" {
			keyTemp := fmt.Sprintf("%v = %v", decodeFields[field.Name], value)
			authHeaderCloudList = append(authHeaderCloudList, keyTemp)
		}
	}

	authHeaderCloud := strings.Join(authHeaderCloudList, ";")

	var authHeaderIMList []string

	fields = reflect.TypeOf(clientConf.Im)
	values = reflect.ValueOf(clientConf.Im)

	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		if decodeFields[field.Name] != "host" {
			value := values.Field(i)
			if value.Interface() != "" {
				keyTemp := fmt.Sprintf("%v = %v", decodeFields[field.Name], value.Interface())
				authHeaderIMList = append(authHeaderIMList, keyTemp)
			}
		}
	}

	authHeaderIM := strings.Join(authHeaderIMList, ";")

	authHeader := authHeaderCloud + "\\n" + authHeaderIM

	//fmt.Printf(authHeader)

	return authHeader
}

// // RefreshToken wraps actions for token refreshing
// func RefreshToken(refreshToken string, clientConf *dodasv1alpha1.Infrastructure) (string, error) {

// 	var token string

// 	clientID := clientConf.Spec.AllowRefresh.ClientID
// 	clientSecret := clientConf.Spec.AllowRefresh.ClientSecret
// 	IAMTokenEndpoint := clientConf.Spec.AllowRefresh.IAMTokenEndpoint

// 	client := &http.Client{
// 		Timeout: 300 * time.Second,
// 	}

// 	req, _ := http.NewRequest("GET", IAMTokenEndpoint, nil)

// 	req.SetBasicAuth(clientID, clientSecret)

// 	req.Header.Set("grant_type", "refresh_token")
// 	req.Header.Set("refresh_token", refreshToken)

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}

// 	body, _ := ioutil.ReadAll(resp.Body)

// 	if resp.StatusCode == 200 {

// 		type accessTokenStruct struct {
// 			AccessToken string `json:"access_token"`
// 		}

// 		var accessTokenJSON accessTokenStruct

// 		err = json.Unmarshal(body, &accessTokenJSON)
// 		if err != nil {
// 			return "", err
// 		}

// 		token = accessTokenJSON.AccessToken

// 	} else {
// 		return "", fmt.Errorf("ERROR: %s", string(body))
// 	}

// 	return token, nil
// }

// MakeRequest function based on inputs
func MakeRequest(request Request) (body []byte, statusCode int, err error) {

	var req *http.Request

	r, err := validateRequest(request)
	if err != nil {
		return nil, -1, fmt.Errorf("Failed to validate request inputs %s", err)
	}

	client := &http.Client{
		Timeout: r.Timeout,
	}

	switch r.RequestType {
	case "POST":
		req, err = http.NewRequest(r.RequestType, r.URL, strings.NewReader(string(r.Content)))
		if err != nil {
			return nil, -1, fmt.Errorf("Failed to create POST http request: %s", err)
		}
	default:
		req, err = http.NewRequest(r.RequestType, r.URL, nil)
		if err != nil {
			return nil, -1, fmt.Errorf("Failed to create %s http request: %s", r.RequestType, err)
		}
	}

	if request.AuthUser != "" && request.AuthPwd != "" {
		req.SetBasicAuth(url.QueryEscape(request.AuthUser), url.QueryEscape(request.AuthPwd))
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	fmt.Println(req.Header.Get("grant_type"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, fmt.Errorf("Remote request failed: %s", err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, fmt.Errorf("Failed to read the response: %s", err)
	}

	return body, resp.StatusCode, nil
}
