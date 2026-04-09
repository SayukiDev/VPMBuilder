package github

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(serverURL string, token string) *Client {
	c := NewClient(token)
	c.resty.SetBaseURL(serverURL)
	return c
}

func TestGetReleases_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/releases" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Accept"); got != "application/vnd.github+json" {
			t.Errorf("unexpected Accept header: %s", got)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[
			{
				"id": 1,
				"tag_name": "v1.0.0",
				"name": "Release 1.0.0",
				"body": "First release",
				"draft": false,
				"prerelease": false,
				"created_at": "2025-01-01T00:00:00Z",
				"published_at": "2025-01-01T00:00:00Z",
				"assets": [
					{
						"name": "package.zip",
						"browser_download_url": "https://example.com/package.zip",
						"size": 1024,
						"content_type": "application/zip"
					}
				]
			},
			{
				"id": 2,
				"tag_name": "v0.9.0",
				"name": "Release 0.9.0",
				"body": "Beta release",
				"draft": false,
				"prerelease": true,
				"created_at": "2024-12-01T00:00:00Z",
				"published_at": "2024-12-01T00:00:00Z",
				"assets": []
			}
		]`))
	}))
	defer srv.Close()

	c := newTestClient(srv.URL, "")
	releases, err := c.GetReleases("owner/repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(releases) != 2 {
		t.Fatalf("expected 2 releases, got %d", len(releases))
	}

	r := releases[0]
	if r.ID != 1 {
		t.Errorf("expected ID 1, got %d", r.ID)
	}
	if r.TagName != "v1.0.0" {
		t.Errorf("expected tag v1.0.0, got %s", r.TagName)
	}
	if r.Name != "Release 1.0.0" {
		t.Errorf("expected name 'Release 1.0.0', got %s", r.Name)
	}
	if r.Draft {
		t.Error("expected draft=false")
	}
	if r.Prerelease {
		t.Error("expected prerelease=false")
	}
	if len(r.Assets) != 1 {
		t.Fatalf("expected 1 asset, got %d", len(r.Assets))
	}
	if r.Assets[0].Name != "package.zip" {
		t.Errorf("expected asset name 'package.zip', got %s", r.Assets[0].Name)
	}
	if r.Assets[0].Size != 1024 {
		t.Errorf("expected asset size 1024, got %d", r.Assets[0].Size)
	}

	if !releases[1].Prerelease {
		t.Error("expected second release prerelease=true")
	}
}

func TestGetReleases_WithToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected 'Bearer test-token', got %q", auth)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(srv.URL, "test-token")
	_, err := c.GetReleases("owner/repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetReleases_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"Not Found"}`))
	}))
	defer srv.Close()

	c := newTestClient(srv.URL, "")
	_, err := c.GetReleases("owner/nonexistent")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

func TestGetReleases_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := newTestClient(srv.URL, "")
	_, err := c.GetReleases("owne/repo")
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

func TestGetReleases_EmptyReleases(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(srv.URL, "")
	releases, err := c.GetReleases("owner/repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(releases) != 0 {
		t.Errorf("expected 0 releases, got %d", len(releases))
	}
}

func TestGetReleases_ConnectionError(t *testing.T) {
	c := newTestClient("http://localhost:1", "")
	_, err := c.GetReleases("owner/repo")
	if err == nil {
		t.Fatal("expected error for connection failure, got nil")
	}
}
