package artsy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiBase = "https://api.artsy.net/api"
const tokenHeader = "X-Xapp-Token"

// GetAccessToken retrieves an access token.
func GetAccessToken(clientID, clientSecret string) (string, error) {
	url := apiBase + "/tokens/xapp_token"
	params := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
	}
	res, err := DoWithParams("POST", "", url, params, nil)
	if err != nil {
		return "", err
	}

	return res["token"].(string), nil
}

func Get(accessToken, endpoint, id string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s/%s", apiBase, endpoint, id)
	return Do("GET", accessToken, url, nil)
}

// Do performs a request
func Do(method, accessToken, url string, headers map[string]string) (map[string]interface{}, error) {
	client := new(http.Client)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if len(accessToken) > 0 {
		req.Header.Set(tokenHeader, accessToken)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	body := map[string]interface{}{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func DoWithParams(method, accessToken, url string, params, headers map[string]string) (map[string]interface{}, error) {
	url += "?"
	i := 0
	for name, value := range params {
		url += name + "=" + value
		if i < len(params)-1 {
			url += "&"
		}
		i++
	}
	return Do(method, accessToken, url, headers)
}
