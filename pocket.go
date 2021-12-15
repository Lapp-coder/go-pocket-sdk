package go_pocket_sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

const (
	host         = "https://getpocket.com/v3"
	authorizeURL = "https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s"

	endpointAdd              = "/add"
	endpointModify           = "/send"
	endpointRetrieving       = "/get"
	endpointRequestToken     = "/oauth/request"
	endpointRequestAuthorize = "/oauth/authorize"

	xErrorHeader     = "X-Error"
	xErrorCodeHeader = "X-Error-Code"

	defaultTimeout = time.Second * 10
)

// Client is a getpocket API client
type Client struct {
	client      *http.Client
	consumerKey string
	redirectURL string
}

// NewClient creates a new client with your application key (to generate a key, create your application here: https://getpocket.com/developer/apps)
func NewClient(consumerKey string) (*Client, error) {
	if consumerKey == "" {
		return nil, ErrEmptyConsumerKey
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

	result, err := c.doHTTP(ctx, endpointRetrieving, req)
	if err != nil {
		return nil, err
	}

	return c.parseItems(result), nil
}

func (c *Client) parseItems(result gjson.Result) []Item {
	var items []Item
	for itemID := range result.Get("list").Map() {
		item := Item{ID: itemID}
		item.fillAllFields(result)
		items = append(items, item)
	}

	return items
}

// Authorize returns the Authorization structure with the access token, username and state obtained from the authorization request
func (c *Client) Authorize(ctx context.Context, requestToken string) (Authorization, error) {
	if requestToken == "" {
		return Authorization{}, ErrEmptyRequestToken
	}

	body := requestAuthorization{
		ConsumerKey: c.consumerKey,
		Code:        requestToken,
	}

	result, err := c.doHTTP(ctx, endpointRequestAuthorize, body)
	if err != nil {
		return Authorization{}, err
	}

	accessToken := result.Get("access_token").String()
	username := result.Get("username").String()
	state := result.Get("state").String()

	if accessToken == "" {
		return Authorization{}, ErrEmptyAccessToken
	}

	return Authorization{
		AccessToken: accessToken,
		Username:    username,
		State:       state,
	}, nil
}

// GetAuthorizationURL returns the url string that is used to grant the user access rights to his Pocket account in your application
func (c Client) GetAuthorizationURL(requestToken string) (string, error) {
	if requestToken == "" {
		return "", ErrEmptyRequestToken
	}

	if c.redirectURL == "" {
		return "", ErrEmptyRedirectURL
	}

	return fmt.Sprintf(authorizeURL, requestToken, c.redirectURL), nil
}

// GetRequestToken returns the request token (code), which will be used later to authenticate the user in your application.
// RedirectURL - where the user will be redirected after authorization (better to specify a link to your application),
// State - metadata string that will be returned at each subsequent authentication response (if you don't need it, specify an empty string).
func (c *Client) GetRequestToken(ctx context.Context, redirectURL string, state string) (string, error) {
	if redirectURL == "" {
		return "", ErrEmptyRedirectURL
	}

	c.redirectURL = redirectURL

	body := requestToken{
		ConsumerKey: c.consumerKey,
		RedirectURL: redirectURL,
		State:       state,
	}

	result, err := c.doHTTP(ctx, endpointRequestToken, body)
	if err != nil {
		return "", err
	}

	token := result.Get("code").String()
	if token == "" {
		return "", ErrEmptyRequestTokenInResponse
	}

	return token, nil
}

func (c *Client) doHTTP(ctx context.Context, endpoint string, body interface{}) (gjson.Result, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("error occurred when marshal the input body: %s", err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, "POST", host+endpoint, bytes.NewBufferString(string(b)))
	if err != nil {
		return gjson.Result{}, fmt.Errorf("error occurred when creating the query: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("error occurred when sending a request to the Pocket server: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return gjson.Result{}, fmt.Errorf("API error %s: %s", resp.Header.Get(xErrorCodeHeader), resp.Header.Get(xErrorHeader))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("error occurred when reading the request body: %s", err.Error())
	}

	result := gjson.Parse(string(respBody))
	if result.String() == "" {
		return gjson.Result{}, ErrFailedToParseInputBody
	}

	return result, nil
}
