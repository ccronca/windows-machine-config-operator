package e2e

import (
	"testing"
	"time"

	"github.com/openshift/windows-machine-config-operator/pkg/controller/windowsmachineconfig/tracker"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// numberOfNodes represent the number of nodes to be dealt with in the test suite.
	// Expose this as a flag.
	numberOfNodes = 1
	// gc is the global context across the test suites.
	gc = globalContext{numberOfNodes: numberOfNodes}
)

// globalContext holds the information that we want to use across the test suites.
// If you want to move item here make sure that
// 1.) It is needed across test suites
// 2.) You're responsible for checking if the field is stale or not. Any field
//     in this struct is not guaranteed to be latest from the apiserver.
type globalContext struct {
	// numberOfNodes to be used for the test suite.
	numberOfNodes int
}

// testContext holds the information related to the individual test suite. This data structure
// should be instantiated by every test suite, so that we can update the test context to be
// passed around to get the information which was created within the test suite. For example,
// if the create test suite creates a Windows Node object, the node object and other related
// information should be easily accessible by other methods within the same test suite.
// Some of the fields we have here can be exposed by via flags to the test suite.
type testContext struct {
	// namespace is the test namespace, we get this from the operator SDK's test framework.
	namespace string
	// osdkTestCtx is the operator sdk framework's test Context
	osdkTestCtx *framework.TestCtx
	// nodes is the list of Windows nodes created by operator
	nodes []v1.Node
	// credentials to be used to access the Windows nodes
	credentials []tracker.Credentials
	// kubeclient is the kube client
	kubeclient kubernetes.Interface
	// tracker is a pointer to the tracker configmap created by operator
	tracker *v1.ConfigMap
	// retryInterval to check for existence of resource in kube api server
	retryInterval time.Duration
	// timeout to terminate checking for the existence of resource in kube apiserver
	timeout time.Duration
}

// NewTestContext returns a new test context to be used by every test.
func NewTestContext(t *testing.T) (*testContext, error) {
	fmwkTestContext := framework.NewTestCtx(t)
	namespace, err := fmwkTestContext.GetNamespace()
	if err != nil {
		return nil, errors.Wrap(err, "test context instantiation failed")
	}
	// number of nodes, retry interval and timeout should come from user-input flags
	return &testContext{osdkTestCtx: fmwkTestContext, kubeclient: framework.Global.KubeClient,
		timeout: time.Minute * 15, retryInterval: time.Second * 5, nodes: make([]v1.Node, 0, gc.numberOfNodes),
		namespace: namespace}, nil
}

// cleanup cleans up the test context
func (tc *testContext) cleanup() {
	tc.osdkTestCtx.Cleanup()
}

func TestMain(m *testing.M) {
	framework.MainEntry(m)
}