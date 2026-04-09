package manifest

type Package struct {
	Name            string            `json:"name"`
	DisplayName     string            `json:"displayName"`
	Version         string            `json:"version"`
	Author          Author            `json:"author"`
	Unity           string            `json:"unity"`
	Description     string            `json:"description"`
	VpmDependencies map[string]string `json:"vpmDependencies"`
	Url             string            `json:"url"`
	LegacyFolders   map[string]string `json:"legacyFolders"`
	LegacyFiles     map[string]string `json:"legacyFiles"`
	LegacyPackages  []string          `json:"legacyPackages"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (p *Package) ToRepoPackage(downloadUrl string) *RepoPackageVersion {
	return &RepoPackageVersion{
		Name:    p.Name,
		Url:     downloadUrl,
		Package: *p,
	}
}
