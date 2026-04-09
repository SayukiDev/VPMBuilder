package builder

import (
	"VPMBuilder/pkg/github"
	"VPMBuilder/pkg/vpm/manifest"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (b *Builder) BuildRepoManifest() error {
	c := github.NewClient("")
	rs, err := c.GetReleases(b.repoUrl)
	if err != nil {
		return err
	}
	if len(rs) == 0 {
		return fmt.Errorf("no releases found for %s", b.repoUrl)
	}
	repo := manifest.NewRepositoryManifest()
	f, err := os.Open(b.repoTemplate)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(repo)
	if err != nil {
		return err
	}
	for _, r := range rs {
		var packagesURL string
		for _, a := range r.Assets {
			if a.Name == "packages.json" {
				packagesURL = a.BrowserDownloadURL
				break
			}
		}
		if packagesURL == "" {
			continue
		}

		body, err := c.GetReleaseAsset(packagesURL)
		if err != nil {
			return fmt.Errorf("get packages.json from release %s: %w", r.TagName, err)
		}
		var pms map[string]*manifest.Package
		if err := json.Unmarshal(body, &pms); err != nil {
			return fmt.Errorf("parse packages.json from release %s: %w", r.TagName, err)
		}
		for _, a := range r.Assets {
			if filepath.Ext(a.Name) != ".zip" {
				continue
			}
			nameS := strings.Split(a.Name, "-")
			pm, ok := pms[nameS[0]]
			if !ok {
				continue
			}
			rpv := pm.ToRepoPackage(a.BrowserDownloadURL)
			if repo.Packages[pm.Name] == nil {
				repo.Packages[pm.Name] = &manifest.RepoPackage{
					Versions: make(map[string]*manifest.RepoPackageVersion),
				}
			}
			repo.Packages[pm.Name].Versions[pm.Version] = rpv
		}
	}
	f2, err := os.OpenFile(filepath.Join(b.output, "repo.json"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f2.Close()
	err = json.NewEncoder(f2).Encode(repo)
	if err != nil {
		return err
	}
	return nil
}

func (b *Builder) MergeRepoManifest(urls []string) error {
	mergedRepo := manifest.NewRepositoryManifest()
	f, err := os.Open(b.repoTemplate)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(mergedRepo)
	if err != nil {
		return err
	}
	for _, url := range urls {
		var repo manifest.RepositoryManifest
		resp, err := b.resty.R().SetResult(&repo).Get(url)
		if err != nil {
			return err
		}
		if resp.IsError() {
			return fmt.Errorf("unexpected status %d for %s: %s", resp.StatusCode(), url, resp.String())
		}
		for k, v := range repo.Packages {
			mergedRepo.Packages[k] = v
		}
	}
	return nil
}
