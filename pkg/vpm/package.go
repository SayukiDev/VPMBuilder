package vpm

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Package struct {
	Name            string            `json:"name" validate:"required"`
	DisplayName     string            `json:"displayName" validate:"required"`
	Version         string            `json:"version" validate:"required"`
	Author          Author            `json:"author" validate:"required"`
	Unity           string            `json:"unity"`
	Description     string            `json:"description"`
	VpmDependencies map[string]string `json:"vpmDependencies,omitempty"`
	Url             string            `json:"url" validate:"required"`
	LegacyFolders   map[string]string `json:"legacyFolders,omitempty"`
	LegacyFiles     map[string]string `json:"legacyFiles,omitempty"`
	LegacyPackages  []string          `json:"legacyPackages,omitempty"`
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
