package redditclient

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type RedditClient struct {
	username     string
	password     string
	clientID     string
	clientSecret string
}

type RedditApiAuth struct {
	AccessToken string `json:"access_token"`
}

func (rc *RedditClient) Auth() (io.ReadCloser, error) {
	body := strings.NewReader(fmt.Sprintf("grant_type=password&username=%v&password=%v", rc.username, rc.password))
	req, err := http.NewRequest(http.MethodPost, "https://www.reddit.com/api/v1/access_token", body)

	if err != nil {
		return nil, err
	}

	basic := fmt.Sprintf("%v:%v", rc.clientID, rc.clientSecret)
	encoded := base64.StdEncoding.EncodeToString([]byte(basic))

	req.Header.Set("authorization", fmt.Sprintf("BASIC %v", encoded))
	req.Header.Set("user-agent", fmt.Sprintf("nfs/0.1 by %v", rc.username))
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)

	return resp.Body, err
}

func (rc *RedditClient) GetSubredditPosts(subredditName, limit string) ([]byte, error) {
	br := make([]byte, 0)
	auth, err := rc.Auth()
	if err != nil {
		return br, err
	}

	token, err := NewRedditApiAuth(auth)
	if err != nil {
		return br, err
	}

	api := fmt.Sprintf("https://oauth.reddit.com/r/%v/new?limit=%v", subredditName, limit)
	req, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		return br, nil
	}

	req.Header.Set("authorization", fmt.Sprintf("Bearer %v", token.AccessToken))
	req.Header.Set("user-agent", fmt.Sprintf("nfs/0.1 by %v", rc.username))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return br, err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)

	return b, err
}

func NewRedditApiAuth(data io.ReadCloser) (RedditApiAuth, error) {
	defer data.Close()

	rdapia := RedditApiAuth{}
	b, err := io.ReadAll(data)

	if err != nil {
		return rdapia, err
	}

	err = json.Unmarshal(b, &rdapia)
	return rdapia, err
}

func NewRedditClient(username, password, clientID, clientSecret string) RedditClient {
	return RedditClient{
		username:     username,
		password:     password,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}
