package envtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/imuxin/ksql/pkg/executor"
	"github.com/imuxin/ksql/pkg/repl"
)

func TestExecuteAndFormat(t *testing.T) {
	env := &envtest.Environment{}
	restConfig, err := env.Start()
	assert.NoError(t, err)

	n := &corev1.Namespace{}
	n.Name = "ksql"

	restConfig.GroupVersion = &corev1.SchemeGroupVersion
	restConfig.APIPath = "api"
	restConfig.NegotiatedSerializer = scheme.Codecs

	restCli, err := rest.RESTClientFor(restConfig)
	assert.NoError(t, err)

	err = restCli.Post().Resource("namespaces").Body(n).Do(context.TODO()).Error()
	assert.NoError(t, err)

	result, err := executor.Execute[unstructured.Unstructured]("select * from namespace", restConfig)
	assert.NoError(t, err)
	fmt.Println(repl.Output(result))

	assert.NoError(t, env.Stop())
}
