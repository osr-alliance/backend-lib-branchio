# Branch.io Go HTTP API Client

## Description

Go implementation for the Branch.io Deep Linking API. Supports create, update and read Branch links.

## API documentation

https://help.branch.io/developers-hub/docs/deep-linking-api


## Installation

```
go get -u github.com/osr-alliance/backend-lib-branchio
```

## Usage example

### Create deep llink

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
	deepLink := &branchiolib.DeepLinkProperties{
		Channel: "backend",
		Feature: "emailVerification",
		Stage:   "Email password sigup",
		Type:     1,
		Tags: []string{"test"},
		Identity: "support@allianceapp.com",
	}
	data := &branchiolib.DeepLinkData{
		DesktopUrl:    "http://allianceapp.com",
		OgDescription: "test",
	}

	customData := map[string]interface{}{
		"groupID":   1,
		"groupName": "myGroup",
	}

	url, err := branchio.CreateLink(ctx, deepLink, data, customData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(url)
```

### Read deep link

```go
    ctx := context.Background()
    linkData, err := branchio.GetLink(ctx, deepLink)
	if err != nil {
		fmt.Println(err)
	}
    fmt.Println(linkData)
```

### Update deep link


Not all the deep link properties can be updated. Read API documentation to check which of them are. 
Trying to update a not updatable field will return error with error `code 0`

```go
    deepLink := &branchiolib.DeepLinkProperties{
		Channel: "twitter",
	}
	customData := map[string]interface{}{
		"groupName": "new group name",
	}
	linkData, err = branchio.UpdateLink(ctx, url, deepLink, nil, customData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(linkData)
```

## Implementation details

According Branch API documentation to create a link there are three types of parameters:

* Deep link properties
* Deep link control parameters
* Custom data which can be passed to the app which are sent as part of the deep link redirecting data

None of this parameters are required, we could create a deep link calling:

```go
url, err := branchio.CreateLink(ctx, nil, nil, nil)
```

The content of the deep link will be created with the default data the API sets in absence of any of those fields and use the Branch account site configuration.

### Parameters explained

In this version it was decided to split the create link parameters according the explanation above, then custom data and deep link redirection data are merged. 
While deep link properties are strong typed, all fields were implemented, deep link control parameters (data) has defined only some fields since there many and and they are out of scope in this release.
Anyway any deep link control parameter can be added in the custom data parameter:

```go
customData := make(map[string]interface{
    $fallback_url = "https://example.com"
}
```

