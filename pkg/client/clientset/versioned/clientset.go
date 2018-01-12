package versioned

import (
	glog "github.com/golang/glog"
	jenkinsv1 "github.com/jenkins-x/jx/pkg/client/clientset/versioned/typed/jx/v1"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	JenkinsV1() jenkinsv1.JenkinsV1Interface
	// Deprecated: please explicitly pick a version if possible.
	Jenkins() jenkinsv1.JenkinsV1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	jenkinsV1 *jenkinsv1.JenkinsV1Client
}

// JenkinsV1 retrieves the JenkinsV1Client
func (c *Clientset) JenkinsV1() jenkinsv1.JenkinsV1Interface {
	return c.jenkinsV1
}

// Deprecated: Jenkins retrieves the default version of JenkinsClient.
// Please explicitly pick a version.
func (c *Clientset) Jenkins() jenkinsv1.JenkinsV1Interface {
	return c.jenkinsV1
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.jenkinsV1, err = jenkinsv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		glog.Errorf("failed to create the DiscoveryClient: %v", err)
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.jenkinsV1 = jenkinsv1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.jenkinsV1 = jenkinsv1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
