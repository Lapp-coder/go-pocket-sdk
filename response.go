package go_pocket_sdk

import (
	"fmt"
	"github.com/tidwall/gjson"
)

type Authorization struct {
	AccessToken string
	Username    string
	State       string
}

type Item struct {
	Id            string
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

func newItem(itemId string) Item {
	return Item{Id: itemId}
}

func (i *Item) fillFields(result gjson.Result) {
	i.ResolvedId = result.Get(fmt.Sprintf("list.%s.resolved_id", i.Id)).String()
	i.GivenUrl = result.Get(fmt.Sprintf("list.%s.given_url", i.Id)).String()
	i.ResolvedUrl = result.Get(fmt.Sprintf("list.%s.resolved_url", i.Id)).String()
	i.GivenTitle = result.Get(fmt.Sprintf("list.%s.given_title", i.Id)).String()
	i.ResolvedTitle = result.Get(fmt.Sprintf("list.%s.resolved_title", i.Id)).String()
	i.Favorite = result.Get(fmt.Sprintf("list.%s.favorite", i.Id)).String()
	i.Status = result.Get(fmt.Sprintf("list.%s.status", i.Id)).String()
	i.Excerpt = result.Get(fmt.Sprintf("list.%s.excerpt", i.Id)).String()
	i.IsArticle = result.Get(fmt.Sprintf("list.%s.is_article", i.Id)).String()
	i.HasImage = result.Get(fmt.Sprintf("list.%s.has_image", i.Id)).String()
	i.HasVideo = result.Get(fmt.Sprintf("list.%s.has_video", i.Id)).String()
	i.WordCount = result.Get(fmt.Sprintf("list.%s.word_count", i.Id)).String()
}
