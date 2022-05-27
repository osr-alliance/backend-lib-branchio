package branchio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseUrl      = "https://api2.branch.io/v1"
	deepLinkPath = "url"
)

type Branchio interface {
	CreateLink(ctx context.Context, deepLink *DeepLinkProperties, data *DeepLinkData, customData map[string]interface{}) (string, error)
	GetLink(ctx context.Context, deeplink string) (string, error)
	UpdateLink(ctx context.Context, deeplink string, deepLinkProperties *DeepLinkProperties, data *DeepLinkData, customData map[string]interface{}) (string, error)
}

type branchio struct {
	branchioClient *http.Client
	key            string
	secret         string
	apiUrl         string
}

func New(branchKey string, branchSecret string) Branchio {
	return &branchio{
		branchioClient: http.DefaultClient,
		key:            branchKey,
		secret:         branchSecret,
		apiUrl:         baseUrl,
	}
}

func (b *branchio) CreateLink(ctx context.Context, deepLink *DeepLinkProperties, data *DeepLinkData, customData map[string]interface{}) (string, error) {
	type UrlMessage struct {
		Url string
	}
	var urlMsg UrlMessage
	dataMap, _ := b.structToMap(data)
	customDataMap := b.mergeMaps(dataMap, customData)
	deepLink.Data = customDataMap
	deepLinkMap, _ := b.structToMap(deepLink)
	deepLinkMap["branch_key"] = b.key

	body, err := b.Post(ctx, deepLinkPath, deepLinkMap)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &urlMsg)
	if err != nil {
		return "", fmt.Errorf("errror parsing response body: %v\n", err)
	}

	return urlMsg.Url, nil
}

func (b *branchio) GetLink(ctx context.Context, deeplink string) (string, error) {
	queryMap := map[string]interface{}{
		"branch_key": b.key,
		"url":        deeplink,
	}

	body, err := b.Get(ctx, deepLinkPath, queryMap)
	if err != nil {
		return "", fmt.Errorf("error Get request: %v", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		return "", fmt.Errorf("error parsing response body: %v\n", err)
	}

	return string(prettyJSON.Bytes()), nil

}

func (b *branchio) UpdateLink(ctx context.Context, deeplink string, deepLinkProperties *DeepLinkProperties, data *DeepLinkData, customData map[string]interface{}) (string, error) {
	queryMap := map[string]interface{}{
		"url": deeplink,
	}
	dataMap, _ := b.structToMap(data)
	customDataMap := b.mergeMaps(dataMap, customData)
	deepLinkProperties.Data = customDataMap
	deepLinkMap, err := b.structToMap(deepLinkProperties)
	if err != nil {
		return "", fmt.Errorf("error parsing deepLinkProperties: %v\n", err)
	}
	deepLinkMap["branch_key"] = b.key
	deepLinkMap["branch_secret"] = b.secret
	fmt.Println("MAP: ", deepLinkMap)
	body, err := b.Put(ctx, deepLinkPath, deepLinkMap, queryMap)
	if err != nil {
		return "", err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t	")
	if err != nil {
		return "", fmt.Errorf("error parsing response body: %v\n", err)
	}

	return string(prettyJSON.Bytes()), nil

}
