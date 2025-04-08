package appscenarios

import (
	"context"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nutanix-cloud-native/nkp-nutanix-product-catalog/apptests/docker"
	"github.com/nutanix-cloud-native/nkp-nutanix-product-catalog/apptests/environment"
	"github.com/nutanix-cloud-native/nkp-nutanix-product-catalog/apptests/flux"
	"github.com/nutanix-cloud-native/nkp-nutanix-product-catalog/apptests/kind"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	genericClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

var (
	env                  *environment.Env
	ctx                  context.Context
	network              *docker.NetworkResource
	k8sClient            genericClient.Client
	restClientV1Pods     rest.Interface
	upgradeAppsRepoPath string
)

var _ = BeforeSuite(func() {
	ctx = context.Background()
	var err error
	network, err = kind.EnsureDockerNetworkExist(ctx, "", false)
	Expect(err).ShouldNot(HaveOccurred())

	env = &environment.Env{
		Network: network,
	}

	// Get the path to upgrade apps repository from the environment variable | (this is where the old version of apps are stored)
	upgradeAppsRepoPath = os.Getenv(upgradeAppsRepoPathEnv)
	if upgradeAppsRepoPath == "" {
		upgradeAppsRepoPath = defaultUpgradeAppsRepoPath
	}
})

func TestApplications(t *testing.T) {
	RegisterFailHandler(Fail)
	suiteConfig, reporterConfig := GinkgoConfiguration()
	RunSpecs(t, "Application Test Suite", suiteConfig, reporterConfig)
}

func SetupKindCluster() error {
	if ctx == nil {
		ctx = context.Background()
	}

	err := env.Provision(ctx)
	if err != nil {
		return err
	}

	k8sClient, err = genericClient.New(env.K8sClient.Config(), genericClient.Options{Scheme: flux.NewScheme()})
	if err != nil {
		return err
	}

	// Get a REST client for making http requests to pods
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Pod",
	}

	httpClient, err := rest.HTTPClientFor(env.K8sClient.Config())
	if err != nil {
		return err
	}

	restClientV1Pods, err = apiutil.RESTClientForGVK(gvk, false, env.K8sClient.Config(), serializer.NewCodecFactory(flux.NewScheme()), httpClient)
	if err != nil {
		return err
	}

	return nil
}
