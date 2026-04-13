package builder

import (
	"VPMBuilder/pkg/github"
	"VPMBuilder/pkg/vpm"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (b *Builder) genPackagesManifest(pms map[string]*vpm.Package) error {
	f, err := os.OpenFile(filepath.Join(b.output, "packages.json"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(pms)
	if err != nil {
		return err
	}
	return nil
}

func (b *Builder) readPackageManifest(path string) (*vpm.Package, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var p vpm.Package
	err = json.NewDecoder(f).Decode(&p)
	if err != nil {
		return nil, err
	}
	err = b.content.Validator.Struct(p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (b *Builder) readRepoManifest(path string) (*vpm.RepositoryManifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var r vpm.RepositoryManifest
	err = json.NewDecoder(f).Decode(&r)
	if err != nil {
		return nil, err
	}
	err = b.content.Validator.Struct(r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (b *Builder) readRepoManifestFromUrl(url string) (*vpm.RepositoryManifest, error) {
	resp, err := b.resty.R().Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("unexpected status %d for %s: %s", resp.StatusCode(), url, resp.String())
	}
	var r vpm.RepositoryManifest
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		return nil, err
	}
	err = b.content.Validator.Struct(r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (b *Builder) getLastPackagesManifest() (map[string]*vpm.Package, error) {
	if b.repoUrl == "" {
		return nil, nil
	}
	c := github.NewClient("")
	rs, err := c.GetReleases(b.repoUrl)
	if err != nil {
		return nil, err
	}
	if len(rs) == 0 {
		return nil, nil
	}
	var url string
	for _, r := range rs {
		for _, a := range r.Assets {
			if a.Name == "packages.json" {
				url = a.BrowserDownloadURL
				break
			}
		}
	}
	if url == "" {
		return nil, nil
	}
	body, err := c.GetReleaseAsset(url)
	if err != nil {
		return nil, err
	}
	pms := make(map[string]*vpm.Package)
	err = json.Unmarshal(body, &pms)
	if err != nil {
		return nil, err
	}
	return pms, nil
}
