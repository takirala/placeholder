package appscenarios

import "time"

const (
	kommanderNamespace = "kommander"
	pollInterval       = 2 * time.Second
	systemClusterCriticalPriority = "system-cluster-critical"
	// Environment variables
	upgradeKappsRepoPathEnv = "UPGRADE_KAPPS_REPO_PATH"
	defaultUpgradeKAppsRepoPath = ".work/upgrade/kommander-applications"
)
