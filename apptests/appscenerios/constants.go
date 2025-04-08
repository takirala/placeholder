package appscenarios

import "time"

const (
	kommanderNamespace = "kommander"
	pollInterval       = 2 * time.Second
	// Environment variables
	upgradeAppsRepoPathEnv = "UPGRADE_KAPPS_REPO_PATH"
	defaultUpgradeAppsRepoPath = ".work/upgrade/applications"
)
