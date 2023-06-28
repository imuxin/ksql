package ext

import (
	"strings"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/ext/abs"
	"github.com/imuxin/ksql/pkg/ext/istio"
	"github.com/imuxin/ksql/pkg/ext/kube"
)

// TODO: refactor this
func NewPlugin(table, namespace string,
	names []string,
	restConfig *rest.Config,
	selector labels.Selector) abs.Plugin {
	switch strings.ToLower(table) {
	case "istio_config":
		return istio.NewPlugin(table, namespace, names, restConfig, selector)
	default:
		return kube.NewPlugin(table, namespace, names, restConfig, selector)
	}
}
