package cmd

import (
	"VPMBuilder/internal/builder"
	"VPMBuilder/internal/log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "merge repo manifest",
	Run:   mergeHandler,
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}

func mergeHandler(_ *cobra.Command, _ []string) {
	b := builder.NewBuilder(config.DownloadRepo, config.OutputPath, config.RepoTemplate, content)
	err := b.MergeRepoManifest(config.RepoUrls)
	if err != nil {
		log.ErrorE("Failed to merge repo manifest", zap.Error(err))
	}
}
