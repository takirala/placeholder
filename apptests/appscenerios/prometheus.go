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

func (pr prometheus) InstallPreviousVersion(ctx context.Context, env *environment.Env) error {
	appPath, err := getPrevVAppsUpgradePath(pr.Name())
	if err != nil {
		return err
	}
	fmt.Printf("==prev version path : %s", appPath)

	err = pr.install(ctx, env, appPath)
	if err != nil {
		return err
	}

	return nil
}

func (pr prometheus) Upgrade(ctx context.Context, env *environment.Env) error {
	appPath, err := absolutePathTo(pr.Name(), "")
	if err != nil {
		return err
	}
	fmt.Printf("==version path to upgrade : %s", appPath)

	err = pr.install(ctx, env, appPath)
	if err != nil {
		return err
	}

	return err
}