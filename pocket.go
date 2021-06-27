package go_pocket_sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/valyala/fastjson"
)

const (
	host         = "https://getpocket.com/v3"
	authorizeUrl = "https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s"

	endpointAdd              = "/add"
	endpointModify           = "/send"
	endpointRetrieving       = "/get"
	endpointRequestToken     = "/oauth/request"
	endpointRequestAuthorize = "/oauth/authorize"

	xErrorHeader     = "X-Error"
	xErrorCodeHeader = "X-Error-Code"

	defaultTimeout = time.Second * 5
)

// Client is a getpocket API client
type Client struct {
	client      *http.Client
	consumerKey string
}

// NewClient creates a new client with your application key (to generate a key, create your application here: https://getpocket.com/developer/apps)
func NewClient(consumerKey string) (*Client, error) {
	if consumerKey == "" {
		return nil, errors.New("empty consumer key")
	}

	return &Client{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		consumerKey: consumerKey,
	}, nil
}

// Add creates a new item in the Pocket list
func (c *Client) Add(ctx context.Context, input AddInput) error {
	req, err := input.generateRequest(c.consumerKey)
	if err != nil {
		return err
	}

	_, err = c.doHTTP(ctx, endpointAdd, req)
	return err
}

// Modify modifies Pocket user data (archives items, adds tags to an item, marks an item as a favorite, etc).
func (c *Client) Modify(ctx context.Context, input ModifyInput) error {
	req, err := input.generateRequest(c.consumerKey)
	if err != nil {
		return err
	}

	_, err = c.doHTTP(ctx, endpointModify, req)
	return err
}

// Retrieving retrieves user data (items) Pocket, such as the item id, which is needed to modify items in the Modify function
func (c *Client) Retrieving(ctx context.Context, input RetrievingInput) ([]Item, error) {
	req, err := input.generateRequest(c.consumerKey)
	if err != nil {
		return nil, err
	}

	values, err := c.doHTTP(ctx, endpointRetrieving, req)
	if err != nil {
		return nil, err
	}

	if values.GetObject("list") == nil {
		return nil, nil
	}

	return c.parseItems(values), nil
}

// Authorize returns the Authorization structure with the access token, user name and state obtained from the authorization request
func (c *Client) Authorize(ctx context.Context, requestToken string) (Authorization, error) {
	if requestToken == "" {
		return Authorization{}, errors.New("empty request token")
	}

	body := requestAuthorization{
		ConsumerKey: c.consumerKey,
		Code:        requestToken,
	}

	values, err := c.doHTTP(ctx, endpointRequestAuthorize, body)
	if err != nil {
		return Authorization{}, err
	}

	accessToken := values.GetStringBytes("access_token")
	username := values.GetStringBytes("username")
	state := values.GetStringBytes("state")
	if string(accessToken) == "" {
		return Authorization{}, errors.New("empty access token in API response")
	}

	return Authorization{
		AccessToken: string(accessToken),
		Username:    string(username),
		State:       string(state),
	}, nil
}

// GetAuthorizationURL returns the url string that is used to grant the user access rights to his Pocket account in your application
func (c *Client) GetAuthorizationURL(requestToken, redirectURL string) (string, error) {
	if requestToken == "" {
		return "", errors.New("empty request token")
	}

	if redirectURL == "" {
		return "", errors.New("empty redirection URL")
	}

	return fmt.Sprintf(authorizeUrl, requestToken, redirectURL), nil
}

// GetRequestToken returns the request token (code), which will be used later to authenticate the user in your application.
// RedirectURL - where the user will be redirected after authorization (better to specify a link to your application),
// State - metadata string that will be returned at each subsequent authentication response (if you don't need it, specify an empty string).
func (c *Client) GetRequestToken(ctx context.Context, redirectURL string, state string) (string, error) {
	if redirectURL == "" {
		return "", errors.New("empty redirect URL")
	}

	body := requestToken{
		ConsumerKey: c.consumerKey,
		RedirectURL: redirectURL,
		State:       state,
	}

	values, err := c.doHTTP(ctx, endpointRequestToken, body)
	if err != nil {
		return "", err
	}

	token := values.GetStringBytes("code")
	if string(token) == "" {
		return "", errors.New("empty request token in API response")
	}

	return string(token), nil
}

func (c *Client) parseItems(values *fastjson.Value) []Item {
	const indexForItemId int = 13

	var (
		items  []Item
		itemId string
		index  int
	)

	newJsonStr := values.GetObject("list").String()
	for index != -1 {
		index = strings.Index(newJsonStr, ":{\"item_id\":\"")
		if index != -1 {
			for i := index + indexForItemId; string(newJsonStr[i]) != "\""; i++ {
				itemId += string(newJsonStr[i])
			}

			items = append(items, createItem(itemId, values))

			itemId = ""
			oldJsonStr := newJsonStr
			newJsonStr = ``

			for i := index + indexForItemId; i != len(oldJsonStr); i++ {
				newJsonStr += string(oldJsonStr[i])
			}
		}
	}

	return items
}

func (c *Client) doHTTP(ctx context.Context, endpoint string, body interface{}) (*fastjson.Value, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("an error occurred when marshal the input body: %s", err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, "POST", host+endpoint, bytes.NewBufferString(string(b)))
	if err != nil {
		return nil, fmt.Errorf("an error occurred when creating the query: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("an error occurred when sending a request to the Pocket server: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %s: %s", resp.Header.Get(xErrorCodeHeader), resp.Header.Get(xErrorHeader))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("an error occurred when reading the request body: %s", err.Error())
	}

	values, err := fastjson.Parse(string(respBody))
	if err != nil {
		return nil, fmt.Errorf("an error occurred while parsing the response body: %s", err.Error())
	}

	return values, nil
}
