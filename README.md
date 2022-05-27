# Branch.io Go HTTP API Client

## Description

Go implementation for the Branch.io API. Supports creating, updating and retrieving Branch deep links.

## Example

```go
package main

import (
	"context"
	"fmt"

	branchiolib "github.com/osr-alliance/backend-lib-branchio"
)

const (
	branch_key    = "KEY"
	branch_secret = "SECRET"
)

func main() {
	ctx := context.Background()
	branchio := branchiolib.New(branch_key, branch_secret)
	deepLink := &branchio.DeepLinkProperties{
		Channel: "backend",
		Feature: "emailVerification",
		Stage:   "Email password sigup",
		Type:     1,
		Tags: []string{"test"},
		Identity: "lucas@allianceapp.com",
	}
	data := &branchio.DeepLinkData{
		DesktopUrl:    "http://allianceapp.com",
		OgDescription: "lucas test",
	}

	customData := map[string]interface{}{
		"groupID":   1,
		"groupName": "lucas test",
	}

	url, err := branchio.CreateLink(ctx, deepLink, data, customData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Returned url: ", url)
```


