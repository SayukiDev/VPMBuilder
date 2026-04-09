package cmd

import (
	"VPMBuilder/internal/builder"
	"VPMBuilder/internal/log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	buildTypePackage = "package"
	buildTypeRepo    = "repo"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build packages or repo manifest",
	Run:   buildHandler,
}
var buildType string

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVarP(&buildType, "type", "t", "package", "build type")
}

func buildHandler(_ *cobra.Command, _ []string) {
	b := builder.NewBuilder(config.DownloadRepo, config.OutputPath, config.RepoTemplate)
	switch buildType {
	case buildTypePackage:
		log.Info("Building packages...", zap.Strings("paths", config.PackagePaths))
		count, err := b.BuildPackages(config.PackagePaths)
		if err != nil {
			log.ErrorE("Failed to build packages", zap.Error(err))
		}
		log.Info("Packages built.", zap.Int("count", count))
	case buildTypeRepo:
		log.Info("Building repo manifest...")
		err := b.BuildRepoManifest()
		if err != nil {
			log.ErrorE("Failed to build repo manifest", zap.Error(err))
		}
		log.Info("Repo manifest built.")
	default:
		log.ErrorE("Invalid build type", zap.String("type", buildType))
	}
}
