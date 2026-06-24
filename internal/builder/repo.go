package builder

import (
	"VPMBuilder/pkg/github"
	"VPMBuilder/pkg/vpm"
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
	repo, err := b.readRepoManifest(b.repoTemplate)
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
		var pms map[string]*vpm.Package
		if err := json.Unmarshal(body, &pms); err != nil {
			return fmt.Errorf("parse packages.json from release %s: %w", r.TagName, err)
		}
		for _, a := range r.Assets {
			if filepath.Ext(a.Name) != ".zip" {
				continue
			}
			nameS := strings.Split(a.Name, "-")
			if len(nameS) > 2 {
				nameS[0] = strings.Join(nameS[:len(nameS)-1], "-")
			}
			pm, ok := pms[nameS[0]]
			if !ok {
				continue
			}
			if pm.ChangelogUrl == "" {
				pm.ChangelogUrl = github.FormatReleaseUrl(b.repoUrl, r.TagName)
			}
			if repo.Packages[pm.Name] == nil {
				repo.Packages[pm.Name] = &vpm.RepoPackageVersions{
					Versions: make(map[string]*vpm.RepoPackage),
				}
			}
			repo.Packages[pm.Name].Versions[pm.Version] = &vpm.RepoPackage{
				Url:     a.BrowserDownloadURL,
				Package: *pm,
			}
		}
	}
	if len(repo.Packages) == 0 {
		return fmt.Errorf("no packages found for %s", b.repoUrl)
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
	mergedRepo, err := b.readRepoManifest(b.repoTemplate)
	if err != nil {
		return err
	}
	for _, url := range urls {
		repo, err := b.readRepoManifestFromUrl(url)
		if err != nil {
			return err
		}
		for k, v := range repo.Packages {
			mergedRepo.Packages[k] = v
		}
	}
	return nil
}
