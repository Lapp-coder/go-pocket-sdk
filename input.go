package go_pocket_sdk

import (
	"errors"
	"strings"
)

// AddInput contains the data needed to create an item in the Pocket list
type AddInput struct {
	AccessToken string
	URL         string
	Title       string
	Tags        []string
	TweetId     string
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

func (ai AddInput) generateRequest(consumerKey string) (requestAdd, error) {
	if ai.AccessToken == "" {
		return requestAdd{}, errors.New("empty access token")
	}

	if ai.URL == "" {
		return requestAdd{}, errors.New("empty URL")
	}

	return requestAdd{
		ConsumerKey: consumerKey,
		AccessToken: ai.AccessToken,
		URL:         ai.URL,
		Title:       ai.Title,
		Tags:        strings.Join(ai.Tags, ", "),
		TweetId:     ai.TweetId,
	}, nil
}

func (mi ModifyInput) generateRequest(consumerKey string) (requestModify, error) {
	if mi.AccessToken == "" {
		return requestModify{}, errors.New("empty access token")
	}

	if len(mi.Actions) == 0 {
		return requestModify{}, errors.New("no actions to modify")
	}

	return requestModify{
		ConsumerKey: consumerKey,
		AccessToken: mi.AccessToken,
		Actions:     mi.Actions,
	}, nil
}

func (ri RetrievingInput) generateRequest(consumerKey string) (requestRetrieving, error) {
	if ri.AccessToken == "" {
		return requestRetrieving{}, errors.New("empty access token")
	}

	return requestRetrieving{
		ConsumerKey: consumerKey,
		AccessToken: ri.AccessToken,
		State:       ri.State,
		Favorite:    ri.Favorite,
		Tag:         ri.Tag,
		ContentType: ri.ContentType,
		Sort:        ri.Sort,
		DetailType:  ri.DetailType,
		Search:      ri.Search,
		Domain:      ri.Domain,
		Since:       ri.Since,
		Count:       ri.Count,
		Offset:      ri.Offset,
	}, nil
}
