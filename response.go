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
	ID            string
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
	return Item{ID: itemId}
}

func (i *Item) fillAllFields(result gjson.Result) {
	i.ResolvedId = result.Get(fmt.Sprintf("list.%s.resolved_id", i.ID)).String()
	i.GivenUrl = result.Get(fmt.Sprintf("list.%s.given_url", i.ID)).String()
	i.ResolvedUrl = result.Get(fmt.Sprintf("list.%s.resolved_url", i.ID)).String()
	i.GivenTitle = result.Get(fmt.Sprintf("list.%s.given_title", i.ID)).String()
	i.ResolvedTitle = result.Get(fmt.Sprintf("list.%s.resolved_title", i.ID)).String()
	i.Favorite = result.Get(fmt.Sprintf("list.%s.favorite", i.ID)).String()
	i.Status = result.Get(fmt.Sprintf("list.%s.status", i.ID)).String()
	i.Excerpt = result.Get(fmt.Sprintf("list.%s.excerpt", i.ID)).String()
	i.IsArticle = result.Get(fmt.Sprintf("list.%s.is_article", i.ID)).String()
	i.HasImage = result.Get(fmt.Sprintf("list.%s.has_image", i.ID)).String()
	i.HasVideo = result.Get(fmt.Sprintf("list.%s.has_video", i.ID)).String()
	i.WordCount = result.Get(fmt.Sprintf("list.%s.word_count", i.ID)).String()
}
