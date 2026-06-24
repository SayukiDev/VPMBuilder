package vpm

type RepoPackage struct {
	Url     string `json:"url" validate:"required"`
	Package `json:",inline"`
}
type RepoPackageVersions struct {
	Versions map[string]*RepoPackage `json:"versions"` // version -> repo package version
}

type RepositoryManifest struct {
	Name     string                          `json:"name"`
	Id       string                          `json:"id"`
	Url      string                          `json:"url"`
	Author   string                          `json:"author"`
	Packages map[string]*RepoPackageVersions `json:"packages"` // name -> repo package
}

func NewRepositoryManifest() *RepositoryManifest {
	return &RepositoryManifest{
		Packages: make(map[string]*RepoPackageVersions),
	}
}
