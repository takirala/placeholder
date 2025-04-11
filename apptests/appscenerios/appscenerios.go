package appscenarios

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nutanix-cloud-native/nkp-nutanix-product-catalog/apptests/environment"
)

type AppScenario interface {
	Install(context.Context, *environment.Env) error // logic implemented by a scenario
	Name() string                                    // scenario name
}

// absolutePathTo returns the absolute path to the given application directory.
func absolutePathTo(application string, appVersion string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// determining the execution path.
	var base string
	_, err = os.Stat(filepath.Join(wd, "services"))
	if os.IsNotExist(err) {
		base = "../.."
	} else {
		base = ""
	}

	dir, err := filepath.Abs(filepath.Join(wd, base, "services", application))
	if err != nil {
		return "", err
	}

	if len(appVersion) > 0 {
		pathToApp := dir + "/" + appVersion
		if !checkIfPathIsValid(pathToApp) {
			return "", fmt.Errorf("no application directory found for app: %s of version: %s", application, appVersion)
		}
		return pathToApp, nil
	}

	// filepath.Glob returns a sorted slice of matching paths
	matches, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", fmt.Errorf(
			"no application directory found for %s in the given path:%s",
			application, dir)
	}

	return matches[0], nil
}

func checkIfPathIsValid(absPathToApp string) bool {
	if _, err := os.Stat(absPathToApp); err == nil {
		return true
	}
	return false
}
