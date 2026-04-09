package manifest

type RepoPackage struct {
	Versions map[string]*RepoPackageVersion `json:"versions"` // version -> repo package version
}

type RepoPackageVersion struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	Package `json:",inline"`
}
type RepositoryManifest struct {
	Name     string                  `json:"name"`
	Id       string                  `json:"id"`
	Url      string                  `json:"url"`
	Author   string                  `json:"author"`
	Packages map[string]*RepoPackage `json:"packages"` // name -> repo package
}

func NewRepositoryManifest() *RepositoryManifest {
	return &RepositoryManifest{
		Packages: make(map[string]*RepoPackage),
	}
}
