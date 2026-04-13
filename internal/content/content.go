package content

import (
	"VPMBuilder/internal/config"

	"github.com/go-playground/validator/v10"
)

type Content struct {
	Cfg       *config.Config
	Validator *validator.Validate
}

func NewContent(cfg *config.Config) *Content {
	v := validator.New()
	return &Content{
		Cfg:       cfg,
		Validator: v,
	}
}
