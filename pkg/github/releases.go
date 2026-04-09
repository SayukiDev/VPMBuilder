package github

import (
	"fmt"
	"slices"

	"github.com/go-resty/resty/v2"
)

func (c *Client) GetReleases(repo string) ([]Release, error) {
	var releases []Release
	var resp *resty.Response
	var err error
	for i := 1; ; i++ {
		temp := make([]Release, 0)
		resp, err = c.resty.R().
			SetResult(&temp).
			SetQueryParams(map[string]string{
				"page": fmt.Sprintf("%d", i),
			}).
			Get("/repos/" + repo + "/releases")
		if err != nil {
			return nil, fmt.Errorf("github: request failed for %s: %w", repo, err)
		}
		if resp.IsError() {
			return nil, fmt.Errorf("github: unexpected status %d for %s: %s", resp.StatusCode(), repo, resp.String())
		}
		releases = append(releases, temp...)
		if len(temp) < 100 {
			break
		}
	}
	slices.Reverse(releases)
	return releases, nil
}

func (c *Client) GetReleaseAsset(url string) ([]byte, error) {
	r, err := c.resty.R().Get(url)
	if err != nil {
		return nil, err
	}
	if r.IsError() {
		return nil, fmt.Errorf("github: unexpected status %d for %s: %s", r.StatusCode(), url, r.String())
	}
	return r.Body(), nil
}
