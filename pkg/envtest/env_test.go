package envtest

import (
	"fmt"
	"testing"

	"github.com/imuxin/ksql/pkg/executor"
	"github.com/imuxin/ksql/pkg/repl"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func TestMain(t *testing.T) {
	env := &envtest.Environment{}
	restConfig, err := env.Start()
	assert.NoError(t, err)

	result, err := executor.Execute[unstructured.Unstructured]("select * from deploy namespace default", restConfig)
	assert.NoError(t, err)
	fmt.Println(repl.Output(result))

	assert.NoError(t, env.Stop())
}
