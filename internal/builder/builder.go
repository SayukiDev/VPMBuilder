package builder

import (
	"VPMBuilder/internal/content"
	"VPMBuilder/internal/log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Builder struct {
	repoUrl      string
	output       string
	repoTemplate string
	resty        *resty.Client
	zap          *zap.Logger
	content      *content.Content
}

func NewBuilder(repoUrl string, output string, repoTemplate string, c *content.Content) *Builder {
	os.MkdirAll(output, 0755)
	return &Builder{
		repoUrl:      repoUrl,
		output:       output,
		repoTemplate: repoTemplate,
		resty:        resty.New().SetTimeout(30 * time.Second),
		zap:          log.SubLogger("builder"),
		content:      c,
	}
}
