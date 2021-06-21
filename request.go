package go_pocket_sdk

type (
	requestToken struct {
		ConsumerKey string `json:"consumer_key"`
		RedirectURL string `json:"redirect_uri"`
		State       string `json:"state,omitempty"`
	}

	requestAuthorization struct {
		ConsumerKey string `json:"consumer_key"`
		Code        string `json:"code"`
	}

	requestAdd struct {
		ConsumerKey string `json:"consumer_key"`
		AccessToken string `json:"access_token"`
		URL         string `json:"url"`
		Title       string `json:"title,omitempty"`
		Tags        string `json:"tags,omitempty"`
		TweetId     string `json:"tweet_id,omitempty"`
	}

	requestModify struct {
		ConsumerKey string   `json:"consumer_key"`
		AccessToken string   `json:"access_token"`
		Actions     []Action `json:"actions"`
	}

	requestRetrieving struct {
		ConsumerKey string `json:"consumer_key"`
		AccessToken string `json:"access_token"`
		State       string `json:"state,omitempty"`
		Favorite    string `json:"favorite,omitempty"`
		Tag         string `json:"tag,omitempty"`
		ContentType string `json:"content_type,omitempty"`
		Sort        string `json:"sort,omitempty"`
		DetailType  string `json:"detail_type,omitempty"`
		Search      string `json:"search,omitempty"`
		Domain      string `json:"domain,omitempty"`
		Since       int64  `json:"since,omitempty"`
		Count       int    `json:"count,omitempty"`
		Offset      int    `json:"offset,omitempty"`
	}
)
