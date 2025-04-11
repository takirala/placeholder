package appscenarios

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/nutanix-cloud-native/nkp-nutanix-product-catalog/apptests/environment"
)

type prometheus struct{}

func (r prometheus) Name() string {
	return "prometheus"
}

func (r prometheus) Install(ctx context.Context, env *environment.Env, appVersion string) error {
	appPath, err := absolutePathTo(r.Name(), appVersion)
	if err != nil {
		return err
	}

	err = r.install(ctx, env, appPath)
	if err != nil {
		return err
	}

	return err
}

func (r prometheus) install(ctx context.Context, env *environment.Env, appPath string) error {

	// apply the kustomization for the source
	sourcesPath := filepath.Join(appPath, "/sources")
	err := env.ApplyKustomizations(ctx, sourcesPath, map[string]string{
		"releaseNamespace": kommanderNamespace,
	})
	if err != nil {
		return err
	}

	// apply the kustomization for the helmrelease
	helmreleasePath := filepath.Join(appPath, "/helmrelease")
	err = env.ApplyKustomizations(ctx, helmreleasePath, map[string]string{
		"releaseNamespace": kommanderNamespace,
	})
	if err != nil {
		return err
	}

	return err
}
