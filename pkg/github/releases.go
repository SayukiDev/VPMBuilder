package github

import "fmt"

func (c *Client) GetReleases(repo string) ([]Release, error) {
	var releases []Release

	resp, err := c.resty.R().
		SetPathParams(map[string]string{
			"repo": repo,
		}).
		SetResult(&releases).
		Get("/repos/{repo}/releases")

	if err != nil {
		return nil, fmt.Errorf("github: request failed for %s: %w", repo, err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("github: unexpected status %d for %s: %s", resp.StatusCode(), repo, resp.String())
	}

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
