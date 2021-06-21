# GetPocket API Golang SDK

[![Release](https://img.shields.io/badge/release-v1.0.0-blue)](https://github.com/Lapp-coder/go-pocket-sdk/releases)

### Основа пакета была сделана на коде из [этого](https://github.com/zhashkevych/go-pocket-sdk) репозитория.

### [English version of the README file](README.md)

***

## Введение:
#### Для создания нового клиента вам потребуется получить consumer key, который выдается при создании вашего приложения на сайте getpocket, а конкретнее [здесь](https://getpocket.com/developer/apps/new). 
#### При создании этого приложения вы можете указать разрешения на использование тех или иных API. 
#### Важно отметить, что если вы не укажите какое-либо из этих разрешений, вы будете получать ошибку, пытаясь вызвать его (API метод) из кода, даже если все остальное пройдет успешно.

## Пример использования:
```go
package main

import (
	"context"
	"fmt"
	pocket "github.com/Lapp-coder/go-pocket-sdk"
	"log"
	"time"
)

func main() {
	ctx := context.Background()

	client, err := pocket.NewClient("<ваш-consumer-key>")
	if err != nil {
		log.Fatal(err)
	}

	requestToken, err := client.GetRequestToken(ctx, "https://google.com", "")
	if err != nil {
		log.Fatal(err)
	}

	authURL, _ := client.GetAuthorizationURL(requestToken, "https://google.com")
	fmt.Println(authURL)

	// Жду, пока пользователь нажмет на ссылку авторизации и предоставит права приложению.
	// Затем продолжаю выполнение программы
	fmt.Scanln()

	auth, err := client.Authorize(ctx, requestToken)
	if err != nil {
		log.Fatal(err)
	}

	// Добавка нового элемента пользователю
	_ = client.Add(ctx, pocket.AddInput{
		AccessToken: auth.AccessToken,
		URL:         "https://github.com",
	})

	// Получение всех элементов пользователя
	items, _ := client.Retrieving(ctx, pocket.RetrievingInput{
		AccessToken: auth.AccessToken,
		Favorite:    "0",
	})

	for _, v := range items {
		// Модифицирование всех найденных элементов пользователя
		actions := []pocket.Action{
			{Name: pocket.Favorite, ItemId: v.ItemId, Time: time.Now().Unix()},
			{Name: pocket.Archive, ItemId: v.ItemId, Time: time.Now().Unix()},
			{Name: pocket.TagsAdd, ItemId: v.ItemId, Tags: "github.com, github, system-version-control"},
		}

		_ = client.Modify(ctx, pocket.ModifyInput{
			AccessToken: auth.AccessToken,
			Actions:     actions,
		})
	}
}
```
