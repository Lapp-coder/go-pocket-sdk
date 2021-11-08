package go_pocket_sdk

import (
	"strings"
)

// AddInput contains the data needed to create an item in the Pocket list
type AddInput struct {
	AccessToken string
	URL         string
	Title       string
	Tags        []string
	TweetID     string
}

// ModifyInput contains the data needed to modify items in the Pocket list
type ModifyInput struct {
	AccessToken string
	Actions     []Action
}

// RetrievingInput contains the data needed to retrieve items from the Pocket list
type RetrievingInput struct {
	AccessToken string
	State       string
	Favorite    string
	Tag         string
	ContentType string
	Sort        string
	DetailType  string
	Search      string
	Domain      string
	Since       int64
	Count       int
	Offset      int
}

func (i AddInput) generateRequest(consumerKey string) (requestAdd, error) {
	if i.AccessToken == "" {
		return requestAdd{}, ErrEmptyAccessToken
	}

	if i.URL == "" {
		return requestAdd{}, ErrEmptyItemURL
	}

	return requestAdd{
		ConsumerKey: consumerKey,
		AccessToken: i.AccessToken,
		URL:         i.URL,
		Title:       i.Title,
		Tags:        strings.Join(i.Tags, ", "),
		TweetID:     i.TweetID,
	}, nil
}

func (i ModifyInput) generateRequest(consumerKey string) (requestModify, error) {
	if i.AccessToken == "" {
		return requestModify{}, ErrEmptyAccessToken
	}

	if len(i.Actions) == 0 {
		return requestModify{}, ErrNoActions
	}

	return requestModify{
		ConsumerKey: consumerKey,
		AccessToken: i.AccessToken,
		Actions:     i.Actions,
	}, nil
}

func (i RetrievingInput) generateRequest(consumerKey string) (requestRetrieving, error) {
	if i.AccessToken == "" {
		return requestRetrieving{}, ErrEmptyAccessToken
	}

	return requestRetrieving{
		ConsumerKey: consumerKey,
		AccessToken: i.AccessToken,
		State:       i.State,
		Favorite:    i.Favorite,
		Tag:         i.Tag,
		ContentType: i.ContentType,
		Sort:        i.Sort,
		DetailType:  i.DetailType,
		Search:      i.Search,
		Domain:      i.Domain,
		Since:       i.Since,
		Count:       i.Count,
		Offset:      i.Offset,
	}, nil
}
