package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

// GetAccessToken ..
func (clientConf *Conf) GetAccessToken(refreshToken string) (token string, err error) {

	clientID := clientConf.AllowRefresh.ClientID
	clientSecret := clientConf.AllowRefresh.ClientSecret
	IAMTokenEndpoint := clientConf.AllowRefresh.IAMTokenEndpoint

	v := url.Values{}

	v.Set("client_id", clientID)
	v.Set("client_secret", clientSecret)
	v.Set("grant_type", "refresh_token")
	v.Set("refresh_token", refreshToken)

	request := Request{
		URL:         IAMTokenEndpoint,
		RequestType: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		AuthUser: clientID,
		AuthPwd:  clientSecret,
		Content:  []byte(v.Encode()),
	}

	body, statusCode, err := MakeRequest(request)
	if err != nil {
		return "", err
	}

	if statusCode != 200 {
		return "", fmt.Errorf("Error code %d: %s", statusCode, string(body))
	}

	var bodyJSON RefreshTokenStruct

	//fmt.Println(string(body))
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		return "", err
	}

	// TODO: only if the mode is IAM for both cloud and
	clientConf.Cloud.Password = bodyJSON.AccessToken
	clientConf.Im.Token = bodyJSON.AccessToken

	return bodyJSON.AccessToken, nil
}

// GetNewToken ..
func (clientConf Conf) GetNewToken() (updatedConf Conf, err error) {

	tokenBytes, err := ioutil.ReadFile(clientConf.AllowRefresh.RefreshTokenFile)
	if err != nil {
		return Conf{}, fmt.Errorf("Failed to read refresh token, please be sure you did `dodas iam init` command: %s", err)
	}

	accessToken, err := clientConf.GetAccessToken(string(tokenBytes))
	if err != nil {
		return Conf{}, err
	}

	//fmt.Printf("Access token: %s", accessToken)

	// TODO: only if the mode is IAM for both cloud and
	clientConf.Cloud.Password = accessToken
	clientConf.Im.Token = accessToken

	// TODO: dump access token to a file

	updatedConf = clientConf
	return updatedConf, nil
}

// GetRefreshToken ..
func (clientConf *Conf) GetRefreshToken() (RefreshToken string, err error) {

	clientID := clientConf.AllowRefresh.ClientID
	clientSecret := clientConf.AllowRefresh.ClientSecret
	IAMTokenEndpoint := clientConf.AllowRefresh.IAMTokenEndpoint
	accessToken := clientConf.Im.Token

	v := url.Values{}

	v.Set("client_id", clientID)
	v.Set("client_secret", clientSecret)
	v.Set("grant_type", "urn:ietf:params:oauth:grant-type:token-exchange")
	v.Set("subject_token", accessToken)

	request := Request{
		URL:         IAMTokenEndpoint,
		RequestType: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		AuthUser: clientID,
		AuthPwd:  clientSecret,
		Content:  []byte(v.Encode()),
	}

	body, statusCode, err := MakeRequest(request)
	if err != nil {
		return "", err
	}

	if statusCode != 200 {
		return "", fmt.Errorf("Error code %d: %s", statusCode, string(body))
	}

	var bodyJSON RefreshTokenStruct

	//fmt.Println(string(body))
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		return "", err
	}

	return bodyJSON.RefreshToken, nil
}
