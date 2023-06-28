package kube

import (
	"errors"
	"time"

	lop "github.com/samber/lo/parallel"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/ext/abs"
	"github.com/imuxin/ksql/pkg/parser"
	"github.com/imuxin/ksql/pkg/pretty"
	"github.com/imuxin/ksql/pkg/util/jsonpath"
)

var DefaultConfigFlags = genericclioptions.NewConfigFlags(true).
	WithDeprecatedPasswordFlag().
	WithDiscoveryBurst(300).
	WithDiscoveryQPS(50.0)

// static (compile time) check that APIServerPlugin satisfies the `Downloader` interface.
var _ abs.Plugin = &APIServerPlugin{}

type APIServerPlugin struct {
	config    *rest.Config
	table     string
	namespace string
	names     []string
	selector  labels.Selector
}

func NewPlugin(
	table, namespace string,
	names []string,
	restConfig *rest.Config,
	selector labels.Selector) abs.Plugin {
	return APIServerPlugin{
		config:    restConfig,
		table:     table,
		namespace: namespace,
		names:     names,
		selector:  selector,
	}
}

func (d APIServerPlugin) RestConfig() (*rest.Config, error) {
	return d.restClientGetter().ToRESTConfig()
}

func (d APIServerPlugin) AllNamespace() bool {
	return d.namespace == ""
}

func (d APIServerPlugin) ResourceTypeOrNameArgs() []string {
	return append([]string{d.table}, d.names...)
}

func (d APIServerPlugin) restClientGetter() resource.RESTClientGetter {
	var wrapper = func(c *rest.Config) *rest.Config {
		r := c
		if d.config != nil {
			r = d.config
		}

		if r.Timeout == 0 {
			r.Timeout = time.Second * 3
		}

		return r
	}

	return DefaultConfigFlags.WithWrapConfigFn(wrapper)
}

func (d APIServerPlugin) Download() ([]abs.Object, error) {
	if d.AllNamespace() && len(d.names) > 1 {
		return nil, errors.New("NAMESPACE required when name is provided")
	}
	r := resource.NewBuilder(d.restClientGetter()).
		Unstructured().
		NamespaceParam(d.namespace).DefaultNamespace().AllNamespaces(d.AllNamespace()).
		LabelSelectorParam(d.selector.String()).
		// FieldSelectorParam(o.FieldSelector).
		// Subresource(o.Subresource).
		RequestChunksOf(0).
		ResourceTypeOrNameArgs(true, d.ResourceTypeOrNameArgs()...).
		// ContinueOnError().
		Latest().
		Flatten().
		// TransformRequests(o.transformRequests).
		Do()
	if r.Err() != nil {
		return nil, r.Err()
	}
	infos, _ := r.Infos()

	return lop.Map(infos, func(item *resource.Info, index int) abs.Object {
		return item.Object.(*unstructured.Unstructured).Object
	}), nil
}

func (d APIServerPlugin) gvk(obj abs.Object) schema.GroupVersionKind {
	apiVersion, _ := jsonpath.Find(obj, "{ .apiVersion }")
	kind, _ := jsonpath.Find(obj, "{ .kind }")
	return schema.FromAPIVersionAndKind(apiVersion, kind)
}

func (d APIServerPlugin) restClientFor(obj abs.Object) (*rest.RESTClient, error) {
	cfg, err := d.RestConfig()
	if err != nil {
		return nil, err
	}

	gvk := d.gvk(obj)
	gv := gvk.GroupVersion()

	cfg.ContentConfig = resource.UnstructuredPlusDefaultContentConfig()
	cfg.GroupVersion = &gv
	if len(gv.Group) == 0 {
		cfg.APIPath = "/api"
	} else {
		cfg.APIPath = "/apis"
	}
	return rest.RESTClientFor(cfg)
}

func (d APIServerPlugin) restMappingFor(obj abs.Object) (*meta.RESTMapping, error) {
	mapper, err := d.restClientGetter().ToRESTMapper()
	if err != nil {
		return nil, err
	}

	gvk := d.gvk(obj)
	return mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
}

func (d APIServerPlugin) resourceHelper(obj abs.Object) (*resource.Helper, error) {
	restClient, err := d.restClientFor(obj)
	if err != nil {
		return nil, err
	}
	restMapping, err := d.restMappingFor(obj)
	if err != nil {
		return nil, err
	}

	return resource.NewHelper(
		restClient, restMapping,
	), nil
}

func (d APIServerPlugin) Delete(list []abs.Object) ([]abs.Object, error) {
	var result []abs.Object
	for _, item := range list {
		helper, err := d.resourceHelper(item)
		if err != nil {
			return nil, err
		}
		namespace, _ := jsonpath.Find(item, "{ .metadata.namespace }")
		name, _ := jsonpath.Find(item, "{ .metadata.name }")
		r, err := helper.Delete(namespace, name)
		if err != nil {
			return nil, err
		}
		result = append(result, r.(*unstructured.Unstructured).Object)
	}
	return result, nil
}

func (d APIServerPlugin) Columns(ksql *parser.KSQL) []pretty.PrintColumn {
	return ksql.CompilePrintColumns()
}
