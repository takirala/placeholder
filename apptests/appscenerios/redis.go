package appscenarios

import (
	"context"
	"path/filepath"

	"github.com/nutanix-cloud-native/nkp-nutanix-product-catalog/apptests/environment"
)

type redis struct{}

func (r redis) Name() string {
	return "redis"
}

func (r redis) Install(ctx context.Context, env *environment.Env) error {
	appPath, err := absolutePathTo(r.Name())
	if err != nil {
		return err
	}

	err = r.install(ctx, env, appPath)
	if err != nil {
		return err
	}

	return err
}

func (r redis) install(ctx context.Context, env *environment.Env, appPath string) error {
	// apply defaults config maps first
	defaultKustomization := filepath.Join(appPath, "/defaults")
	err := env.ApplyKustomizations(ctx, defaultKustomization, map[string]string{
		"releaseNamespace": kommanderNamespace,
	})
	if err != nil {
		return err
	}

	// apply the kustomization for the source
	sourcesPath := filepath.Join(appPath, "/sources")
	err = env.ApplyKustomizations(ctx, sourcesPath, map[string]string{
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
