package go_pocket_sdk

import (
	"github.com/valyala/fastjson"
)

type Authorization struct {
	AccessToken string
	Username    string
	State       string
}

type Item struct {
	ItemId        string
	ResolvedId    string
	GivenUrl      string
	ResolvedUrl   string
	GivenTitle    string
	ResolvedTitle string
	Favorite      string
	Status        string
	Excerpt       string
	IsArticle     string
	HasImage      string
	HasVideo      string
	WordCount     string
}

func createItem(itemId string, values *fastjson.Value) Item {
	return Item{
		ItemId:        itemId,
		ResolvedId:    string(values.GetStringBytes("list", itemId, "resolved_id")),
		GivenUrl:      string(values.GetStringBytes("list", itemId, "given_url")),
		ResolvedUrl:   string(values.GetStringBytes("list", itemId, "resolved_url")),
		GivenTitle:    string(values.GetStringBytes("list", itemId, "given_title")),
		ResolvedTitle: string(values.GetStringBytes("list", itemId, "resolved_title")),
		Favorite:      string(values.GetStringBytes("list", itemId, "favorite")),
		Status:        string(values.GetStringBytes("list", itemId, "status")),
		Excerpt:       string(values.GetStringBytes("list", itemId, "excerpt")),
		IsArticle:     string(values.GetStringBytes("list", itemId, "is_article")),
		HasImage:      string(values.GetStringBytes("list", itemId, "has_image")),
		HasVideo:      string(values.GetStringBytes("list", itemId, "has_video")),
		WordCount:     string(values.GetStringBytes("list", itemId, "word_count")),
	}
}
