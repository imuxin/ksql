package kube

import (
	lop "github.com/samber/lo/parallel"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/ext"
)

// static (compile time) check that APIServerDownloader satisfies the `Downloader` interface.
var _ ext.Downloader = &APIServerDownloader{}

type APIServerDownloader struct {
	RestConfig *rest.Config
	Database   string
	Table      string
	Namespace  string
	Names      []string
	Selector   labels.Selector
}

func (d APIServerDownloader) AllNamespace() bool {
	return d.Namespace == ""
}

func (d APIServerDownloader) ResourceTypeOrNameArgs() []string {
	return append([]string{d.Table}, d.Names...)
}

func (d APIServerDownloader) restClientGetter() resource.RESTClientGetter {
	var wrapper = func(c *rest.Config) *rest.Config {
		if d.RestConfig != nil {
			return d.RestConfig
		}
		return c
	}

	return genericclioptions.NewConfigFlags(true).
		WithDeprecatedPasswordFlag().
		WithDiscoveryBurst(300).
		WithDiscoveryQPS(50.0).
		WithWrapConfigFn(wrapper)
}

func (d APIServerDownloader) Download() ([]ext.Object, error) {
	r := resource.NewBuilder(d.restClientGetter()).
		Unstructured().
		NamespaceParam(d.Namespace).DefaultNamespace().AllNamespaces(d.AllNamespace()).
		LabelSelectorParam(d.Selector.String()).
		// FieldSelectorParam(o.FieldSelector).
		// Subresource(o.Subresource).
		RequestChunksOf(0).
		ResourceTypeOrNameArgs(true, d.ResourceTypeOrNameArgs()...).
		// ContinueOnError().
		Latest().
		Flatten().
		// TransformRequests(o.transformRequests).
		Do()
	infos, _ := r.Infos()

	return lop.Map(infos, func(item *resource.Info, index int) ext.Object {
		return item.Object.(*unstructured.Unstructured).Object
	}), nil
}
