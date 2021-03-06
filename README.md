# GetPocket API Golang SDK

[![Release](https://img.shields.io/badge/release-v1.0.5-blue)](https://github.com/Lapp-coder/go-pocket-sdk/releases)

### The basis of the package was made on code from [this](https://github.com/zhashkevych/go-pocket-sdk) repository.

### [Русская версия README файла](README_RU.md)

***

## Installation
```go get -u github.com/Lapp-coder/go-pocket-sdk```

## Introduction
#### To create a new client, you need to get the consumer key that you get when you create your application on the getpocket website, specifically [here](https://getpocket.com/developer/apps/new)
#### When you create this application, you can specify the permissions to use of certain APIs.
#### It is important to note that if you do not specify any of these permissions, you will get an error when trying to call API method from the code, even if everything else goes well.


## Example usage:
```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pocket "github.com/Lapp-coder/go-pocket-sdk"
)

func main() {
	ctx := context.Background()

	client, err := pocket.NewClient("<your-consumer-key>")
	if err != nil {
		log.Fatal(err)
	}

	requestToken, err := client.GetRequestToken(ctx, "https://google.com", "")
	if err != nil {
		log.Fatal(err)
	}

	authURL, _ := client.GetAuthorizationURL(requestToken)
	fmt.Println(authURL)

	// Waiting for the user to follow the authorization link and grant rights to the application
	fmt.Scanln()

	auth, err := client.Authorize(ctx, requestToken)
	if err != nil {
		log.Fatal(err)
	}

	// Adding a new element
	_ = client.Add(ctx, pocket.AddInput{
		AccessToken: auth.AccessToken,
		URL:         "https://github.com",
	})

	// Getting all user items
	items, _ := client.Retrieving(ctx, pocket.RetrievingInput{
		AccessToken: auth.AccessToken,
		Favorite:    "0",
	})

	for _, item := range items {
		// Modifying all found user elements
		actions := []pocket.Action{
			{Name: pocket.ActionFavorite, ItemID: item.ID, Time: time.Now().Unix()},
			{Name: pocket.ActionArchive, ItemID: item.ID, Time: time.Now().Unix()},
			{Name: pocket.ActionTagsAdd, ItemID: item.ID, Tags: "github.com, github, system-version-control"},
		}

		_ = client.Modify(ctx, pocket.ModifyInput{
			AccessToken: auth.AccessToken,
			Actions:     actions,
		})
	}
}
```
