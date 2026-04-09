package github

import "github.com/go-resty/resty/v2"

const baseURL = "https://api.github.com"

type Client struct {
	resty *resty.Client
}

func NewClient(token string) *Client {
	r := resty.New().
		SetBaseURL(baseURL).
		SetHeader("Accept", "application/vnd.github+json")

	if token != "" {
		r.SetAuthToken(token)
	}

	return &Client{resty: r}
}
