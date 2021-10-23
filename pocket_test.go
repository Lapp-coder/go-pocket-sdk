package go_pocket_sdk

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (s roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return s(r)
}

func newClient(t *testing.T, statusCode int, path, responseBody string) *Client {
	return &Client{
		client: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, path, r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)

				return &http.Response{
					StatusCode: statusCode,
					Body:       ioutil.NopCloser(strings.NewReader(responseBody)),
				}, nil
			}),
		},
	}
}

func TestClient_Add(t *testing.T) {
	type args struct {
		ctx      context.Context
		addInput AddInput
	}

	testTable := []struct {
		name                 string
		input                args
		expectedStatusCode   int
		expectedResponseBody string
		expectedErrorMessage string
		wantErr              bool
	}{
		{
			name: "OK_AllFields",
			input: args{
				ctx: context.Background(),
				addInput: AddInput{
					AccessToken: "access-token",
					URL:         "https://github.com",
					Title:       "title",
					Tags:        []string{"github"},
					TweetId:     "1",
				},
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":1}`,
			wantErr:              false,
		},
		{
			name: "OK_WithoutTweetId",
			input: args{
				ctx: context.Background(),
				addInput: AddInput{
					AccessToken: "access-token",
					URL:         "https://github.com",
					Title:       "title",
					Tags:        []string{"github"},
				},
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":1}`,
			wantErr:              false,
		},
		{
			name: "OK_WithoutTagsAndTweetId",
			input: args{
				ctx: context.Background(),
				addInput: AddInput{
					AccessToken: "access-token",
					URL:         "https://github.com",
					Title:       "title",
				},
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":1}`,
			wantErr:              false,
		},
		{
			name: "OK_WithoutTitleAndTagsAndTweetId",
			input: args{
				ctx: context.Background(),
				addInput: AddInput{
					AccessToken: "access-token",
					URL:         "https://github.com",
				},
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":1}`,
			wantErr:              false,
		},
		{
			name: "Empty access token",
			input: args{
				ctx: context.Background(),
				addInput: AddInput{
					AccessToken: "",
					URL:         "https://github.com",
				},
			},
			expectedResponseBody: `{"status":0}`,
			expectedErrorMessage: "empty access token",
			wantErr:              true,
		},
		{
			name: "Empty URL",
			input: args{
				ctx: context.Background(),
				addInput: AddInput{
					AccessToken: "access-token",
					URL:         "",
				},
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":0}`,
			expectedErrorMessage: "empty URL",
			wantErr:              true,
		},
		{
			name: "Non-2XX response",
			input: args{
				ctx: context.Background(),
				addInput: AddInput{
					AccessToken: "access-token",
					URL:         "https://github.com",
					Title:       "title",
					Tags:        []string{"github"},
					TweetId:     "1",
				},
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":0}`,
			expectedErrorMessage: "API error : ",
			wantErr:              true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			client := newClient(t, tc.expectedStatusCode, "/v3/add", tc.expectedResponseBody)

			err := client.Add(tc.input.ctx, tc.input.addInput)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErrorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_Modify(t *testing.T) {
	type args struct {
		ctx         context.Context
		modifyInput ModifyInput
	}

	testTable := []struct {
		name                 string
		input                args
		expectedStatusCode   int
		expectedResponseBody string
		expectedErrorMessage string
		wantErr              bool
	}{
		{
			name: "OK",
			input: args{
				ctx: context.Background(),
				modifyInput: ModifyInput{
					AccessToken: "access-token",
					Actions: []Action{
						{Name: ActionAdd, ItemID: "987", Url: "https://github.com", Title: "Github"},
						{Name: ActionArchive, ItemID: "654", Time: time.Now().Unix()},
						{Name: ActionFavorite, ItemID: "321", Time: time.Now().Unix()},
					},
				},
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":1}`,
			wantErr:              false,
		},
		{
			name: "Empty access token",
			input: args{
				ctx: context.Background(),
				modifyInput: ModifyInput{
					AccessToken: "",
					Actions: []Action{
						{Name: ActionAdd, ItemID: "987", Url: "https://github.com", Title: "Github"},
						{Name: ActionArchive, ItemID: "654", Time: time.Now().Unix()},
						{Name: ActionFavorite, ItemID: "321", Time: time.Now().Unix()},
					},
				},
			},
			expectedResponseBody: `{"status":0}`,
			expectedErrorMessage: "empty access token",
			wantErr:              true,
		},
		{
			name: "Empty array actions",
			input: args{
				ctx: context.Background(),
				modifyInput: ModifyInput{
					AccessToken: "access-token",
					Actions:     []Action{},
				},
			},
			expectedResponseBody: `{"status":0}`,
			expectedErrorMessage: "no actions to modify",
			wantErr:              true,
		},
		{
			name: "Non-2XX response",
			input: args{
				ctx: context.Background(),
				modifyInput: ModifyInput{
					AccessToken: "access-token",
					Actions: []Action{
						{Name: ActionAdd, ItemID: "987", Url: "https://github.com", Title: "Github"},
						{Name: ActionArchive, ItemID: "654", Time: time.Now().Unix()},
						{Name: ActionFavorite, ItemID: "321", Time: time.Now().Unix()},
					},
				},
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":0}`,
			expectedErrorMessage: "API error : ",
			wantErr:              true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			client := newClient(t, tc.expectedStatusCode, "/v3/send", tc.expectedResponseBody)

			err := client.Modify(tc.input.ctx, tc.input.modifyInput)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErrorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_Retrieving(t *testing.T) {
	type args struct {
		ctx             context.Context
		retrievingInput RetrievingInput
	}

	testTable := []struct {
		name                 string
		input                args
		expectedStatusCode   int
		expectedResponseBody string
		expectedItems        []Item
		expectedErrorMessage string
		wantErr              bool
	}{
		{
			name: "OK",
			input: args{
				ctx: context.Background(),
				retrievingInput: RetrievingInput{
					AccessToken: "access-token",
					Favorite:    "0",
				},
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":1,"list":{"229279689":{"item_id":"229279689","resolved_id":"229279689","given_url":"http:\/\/www.grantland.com\/blog\/the-triangle\/post\/_\/id\/38347\/ryder-cup-preview","given_title":"The Massive Ryder Cup Preview - The Triangle Blog - Grantland","favorite":"0","status":"0","resolved_title":"The Massive Ryder Cup Preview","resolved_url":"http:\/\/www.grantland.com\/blog\/the-triangle\/post\/_\/id\/38347\/ryder-cup-preview","excerpt":"The list of things I love about the Ryder Cup is so long that it could fill a (tedious) novel, and golf fans can probably guess most of them.","is_article":"1","has_video":"1","has_image":"1","word_count":"3197"}}}`,
			expectedItems: []Item{
				{
					ID:            "229279689",
					ResolvedId:    "229279689",
					GivenUrl:      `http://www.grantland.com/blog/the-triangle/post/_/id/38347/ryder-cup-preview`,
					GivenTitle:    `The Massive Ryder Cup Preview - The Triangle Blog - Grantland`,
					ResolvedUrl:   `http://www.grantland.com/blog/the-triangle/post/_/id/38347/ryder-cup-preview`,
					ResolvedTitle: `The Massive Ryder Cup Preview`,
					Favorite:      "0",
					Status:        "0",
					Excerpt:       `The list of things I love about the Ryder Cup is so long that it could fill a (tedious) novel, and golf fans can probably guess most of them.`,
					IsArticle:     "1",
					HasImage:      "1",
					HasVideo:      "1",
					WordCount:     "3197",
				},
			},
			wantErr: false,
		},
		{
			name: "Empty access token",
			input: args{
				ctx: context.Background(),
				retrievingInput: RetrievingInput{
					AccessToken: "",
					Favorite:    "0",
				},
			},
			expectedErrorMessage: "empty access token",
			wantErr:              true,
		},
		{
			name: "Non-2XX response",
			input: args{
				ctx: context.Background(),
				retrievingInput: RetrievingInput{
					AccessToken: "access-token",
					Favorite:    "0",
				},
			},
			expectedStatusCode:   400,
			expectedErrorMessage: "API error : ",
			wantErr:              true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			client := newClient(t, tc.expectedStatusCode, "/v3/get", tc.expectedResponseBody)

			got, err := client.Retrieving(tc.input.ctx, tc.input.retrievingInput)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErrorMessage, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedItems, got)
			}
		})
	}
}

func TestClient_Authorize(t *testing.T) {
	type args struct {
		ctx          context.Context
		requestToken string
	}

	testTable := []struct {
		name                  string
		input                 args
		expectedStatusCode    int
		expectedResponseBody  string
		expectedAuthorization Authorization
		expectedErrorMessage  string
		wantErr               bool
	}{
		{
			name: "OK_AllFields",
			input: args{
				ctx:          context.Background(),
				requestToken: "request-token",
			},
			expectedStatusCode:    200,
			expectedResponseBody:  `{"access_token":"access-token","username":"pocket-user","state":"testing"}`,
			expectedAuthorization: Authorization{AccessToken: "access-token", Username: "pocket-user", State: "testing"},
			wantErr:               false,
		},
		{
			name: "OK_WithoutState",
			input: args{
				ctx:          context.Background(),
				requestToken: "request-token",
			},
			expectedStatusCode:    200,
			expectedResponseBody:  `{"access_token":"access-token","username":"pocket-user"}`,
			expectedAuthorization: Authorization{AccessToken: "access-token", Username: "pocket-user"},
			wantErr:               false,
		},
		{
			name: "Empty request token",
			input: args{
				ctx:          context.Background(),
				requestToken: "",
			},
			expectedErrorMessage: "empty request token",
			wantErr:              true,
		},
		{
			name: "Non-2XX response",
			input: args{
				ctx:          context.Background(),
				requestToken: "request-token",
			},
			expectedStatusCode:   400,
			expectedErrorMessage: "API error : ",
			wantErr:              true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			client := newClient(t, tc.expectedStatusCode, "/v3/oauth/authorize", tc.expectedResponseBody)

			got, err := client.Authorize(tc.input.ctx, tc.input.requestToken)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErrorMessage, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedAuthorization, got)
			}
		})
	}
}

func TestClient_GetAuthorizationURL(t *testing.T) {
	type args struct {
		ctx          context.Context
		requestToken string
		redirectURL  string
	}

	expectedAuthorizationURL := func(input args) string {
		return fmt.Sprintf("https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s", input.requestToken, input.redirectURL)
	}

	testTable := []struct {
		name                 string
		input                args
		expectedErrorMessage string
		wantErr              bool
	}{
		{
			name: "OK",
			input: args{
				ctx:          context.Background(),
				requestToken: "request-token",
				redirectURL:  "http://localhost",
			},
			wantErr: false,
		},
		{
			name: "Empty request token",
			input: args{
				ctx:          context.Background(),
				requestToken: "",
				redirectURL:  "http://localhost",
			},
			expectedErrorMessage: "empty request token",
			wantErr:              true,
		},
		{
			name: "Empty redirect URL",
			input: args{
				ctx:          context.Background(),
				requestToken: "request-token",
				redirectURL:  "",
			},
			expectedErrorMessage: "empty redirection URL",
			wantErr:              true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			client := &Client{redirectURL: tc.input.redirectURL}

			got, err := client.GetAuthorizationURL(tc.input.requestToken)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErrorMessage, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedAuthorizationURL(tc.input), got)
			}
		})
	}
}

func TestClient_GetRequestToken(t *testing.T) {
	type args struct {
		ctx         context.Context
		state       string
		redirectURL string
	}

	testTable := []struct {
		name                 string
		input                args
		expectedStatusCode   int
		expectedResponseBody string
		expectedRequestToken string
		expectedErrorMessage string
		wantErr              bool
	}{
		{
			name: "OK_WithState",
			input: args{
				ctx:         context.Background(),
				state:       "testing",
				redirectURL: "http://localhost",
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"code":"request-token","state":"testing"}`,
			expectedRequestToken: "request-token",
			wantErr:              false,
		},
		{
			name: "OK_WithoutState",
			input: args{
				ctx:         context.Background(),
				state:       "",
				redirectURL: "http://localhost",
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"code":"request-token","state":null}`,
			expectedRequestToken: "request-token",
			wantErr:              false,
		},
		{
			name: "Empty redirect URL",
			input: args{
				ctx:         context.Background(),
				redirectURL: "",
			},
			expectedErrorMessage: "empty redirect URL",
			wantErr:              true,
		},
		{
			name: "Empty response code",
			input: args{
				ctx:         context.Background(),
				redirectURL: "http://localhost",
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"code":""}`,
			expectedErrorMessage: "empty request token in API response",
			wantErr:              true,
		},
		{
			name: "Non-2XX response",
			input: args{
				ctx:         context.Background(),
				redirectURL: "http://localhost",
			},
			expectedStatusCode:   400,
			expectedErrorMessage: "API error : ",
			wantErr:              true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			client := newClient(t, tc.expectedStatusCode, "/v3/oauth/request", tc.expectedResponseBody)

			got, err := client.GetRequestToken(tc.input.ctx, tc.input.redirectURL, tc.input.state)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErrorMessage, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRequestToken, got)
			}
		})
	}
}
