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

func (b *branchio) CreateLink(ctx context.Context, deepLinkProperties *DeepLinkProperties, data *DeepLinkData, customData map[string]interface{}) (string, error) {
	type UrlMessage struct {
		Url string
	}
	var urlMsg UrlMessage

	deepLinkParams, err := b.buildParameters(deepLinkProperties, data, customData)
	if err != nil {
		return "", err
	}

	deepLinkParams["branch_key"] = b.key

	body, err := b.Post(ctx, deepLinkPath, deepLinkParams)
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

	deepLinkParams, err := b.buildParameters(deepLinkProperties, data, customData)
	if err != nil {
		return "", err
	}

	deepLinkParams["branch_key"] = b.key
	deepLinkParams["branch_secret"] = b.secret

	body, err := b.Put(ctx, deepLinkPath, deepLinkParams, queryMap)
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
