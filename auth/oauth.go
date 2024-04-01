package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

func Oauth2Request(scope, refreshToken string) (string, string, error) {
	tenantID := viper.Get("tenant")
	oauth := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	data := url.Values{
		"client_id":     {"AAa-aaaa-bbbb-bbb-78787346"}, // TODO replace with your client ID
		"scope":         {scope},
		"grant_type":    {"refresh_token"},
		"client_info":   {"1"},
		"claims":        {"{\"access_token\":{\"xms_cc\":{\"values\":[\"CP1\"]}}}"},
		"refresh_token": {refreshToken},
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", oauth, strings.NewReader(data.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://teams.microsoft.com")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	bodyB, _ := io.ReadAll(resp.Body)
	var json_data map[string]interface{}
	json.Unmarshal(bodyB, &json_data)

	access_token, ok := json_data["access_token"].(string)
	if !ok {
		//? In case we get an error let's check for description
		error_description, ok := json_data["error_description"].(string)
		if !ok {
			log.Fatal("Uknown error with parsing initial token response")
		}
		log.Fatal(error_description)
	}
	new_refresh_token, ok := json_data["refresh_token"].(string)
	if !ok {
		//? In case we get an error let's check for description
		error_description, ok := json_data["error_description"].(string)
		if !ok {
			log.Fatal("Uknown error with parsing initial token response")
		}
		log.Fatal(error_description)
	}
	return access_token, new_refresh_token, nil
}

func GetTeamsToken(tenantID, refresh string) (access_token, refresh_token string) {
	access_token, refresh_token, _ = Oauth2Request(tenantID, "https://chatsvcagg.teams.microsoft.com/.default openid profile offline_access", refresh)
	return
}

func GetSkypeToken(tenantID, refresh string) string {
	access_token, _, _ := Oauth2Request(tenantID, "https://api.spaces.skype.com/.default openid profile offline_access", refresh)

	skype_authz := "https://teams.microsoft.com/api/authsvc/v1.0/authz"

	client := &http.Client{}
	req, _ := http.NewRequest("POST", skype_authz, nil)

	req.Header.Set("content-length", "0")
	req.Header.Set("authorization", "Bearer "+access_token)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("ERROR: SKYPE http request")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyB, _ := io.ReadAll(resp.Body)

	var json_data map[string]interface{}
	json.Unmarshal(bodyB, &json_data)

	tokens_temp, ok := json_data["tokens"].(map[string]interface{})
	if !ok {
		log.Println(resp.StatusCode)
		log.Println(string(bodyB))
		log.Fatal("Error no tokens!")
	}
	skype_token := tokens_temp["skypeToken"].(string)
	return skype_token
}
