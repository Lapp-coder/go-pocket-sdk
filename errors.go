package go_pocket_sdk

import (
	"fmt"
)

var (
	ErrEmptyConsumerKey            = fmt.Errorf("empty consumer key")
	ErrEmptyRequestToken           = fmt.Errorf("empty request token")
	ErrEmptyAccessToken            = fmt.Errorf("empty access token")
	ErrEmptyRedirectURL            = fmt.Errorf("empty redirect URL")
	ErrEmptyRequestTokenInResponse = fmt.Errorf("empty request token in API response")
	ErrEmptyItemURL                = fmt.Errorf("empty URL for add item")
	ErrNoActions                   = fmt.Errorf("no actions to modify items")
	ErrFailedToParseInputBody      = fmt.Errorf("failed to parse input body")
)
