package builder

import (
	"VPMBuilder/pkg/vpm/manifest"
	"VPMBuilder/pkg/zip"
	"context"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/sourcegraph/conc/pool"
	"go.uber.org/zap"
)

func (b *Builder) BuildPackages(path []string) (int, error) {
	m, err := b.getLastPackagesManifest()
	if err != nil {
		return 0, err
	}
	if m == nil {
		m = make(map[string]*manifest.Package)
	}
	pl := pool.New().
		WithContext(context.Background()).
		WithCancelOnError().
		WithFirstError().
		WithMaxGoroutines(runtime.NumCPU() * 2)
	pmChan := make(chan *manifest.Package, len(path))
	for _, p := range path {
		pl.Go(func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return nil
			default:
			}
			b.zap.Debug("Building Package...", zap.String("path", p))
			pm, err := readPackageManifest(filepath.Join(p, "package.json"))
			if err != nil {
				return fmt.Errorf("read package.json failed: %v", err)
			}
			if p, ok := m[pm.Name]; ok {
				if p.Version == pm.Version {
					return nil
				}
			}
			err = b.buildPackage(p, pm)
			if err != nil {
				return err
			}
			pmChan <- pm
			b.zap.Debug("Package built successfully", zap.String("path", p))
			return nil
		})
	}
	go func() {
		err = pl.Wait()
		close(pmChan)
	}()
	for pm := range pmChan {
		m[pm.Name] = pm
	}
	if err != nil {
		return 0, err
	}
	b.zap.Debug("Generating packages.json")
	err = b.genPackagesManifest(m)
	if err != nil {
		return 0, err
	}
	b.zap.Debug("Packages.json generated successfully")
	return len(m), nil
}

func (b *Builder) buildPackage(path string, p *manifest.Package) error {
	err := zip.Compress(path, filepath.Join(
		b.output,
		fmt.Sprintf("%s-%s.zip", p.Name, p.Version),
	))
	if err != nil {
		return fmt.Errorf("compress package failed: %v", err)
	}
	return nil
}
