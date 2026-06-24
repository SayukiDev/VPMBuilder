package vpm

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Package struct {
	Name             string            `json:"name" validate:"required"`
	DisplayName      string            `json:"displayName" validate:"required"`
	Version          string            `json:"version" validate:"required"`
	Author           Author            `json:"author" validate:"required"`
	Unity            string            `json:"unity"`
	Keywords         []string          `json:"keywords,omitempty"`
	Description      string            `json:"description"`
	VpmDependencies  map[string]string `json:"vpmDependencies,omitempty"`
	LegacyFolders    map[string]string `json:"legacyFolders,omitempty"`
	LegacyFiles      map[string]string `json:"legacyFiles,omitempty"`
	LegacyPackages   []string          `json:"legacyPackages,omitempty"`
	HideInEditor     bool              `json:"hideInEditor,omitempty"`
	ChangelogUrl     string            `json:"changelogUrl,omitempty"`
	DocumentationUrl string            `json:"documentationUrl,omitempty"`
	Samples          any               `json:"samples,omitempty"`
	Type             string            `json:"type,omitempty"`
	UnityRelease     string            `json:"unityRelease,omitempty"`
	ZipSHA256        string            `json:"zipSHA256"`
}

// Validate validates the Package struct
func (p *Package) Validate() error {
	return validate.Struct(p)
}

type Author struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Url   string `json:"url"`
}
