package runtime

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
)

var defaultConfigFlags = genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)

type Downloader interface {
	Download() (*unstructured.UnstructuredList, error)
}

// static (compile time) check that APIServerDownloader satisfies the `Downloader` interface.
var _ Downloader = &APIServerDownloader{}

type APIServerDownloader struct {
	Database  string
	Table     string
	Namespace string
	Names     []string
	Selector  labels.Selector
}

func (d APIServerDownloader) ResourceTypeOrNameArgs() []string {
	return append([]string{d.Table}, d.Names...)
}

func (d APIServerDownloader) Download() (*unstructured.UnstructuredList, error) {
	r := resource.NewBuilder(defaultConfigFlags).
		Unstructured().
		NamespaceParam(d.Namespace).DefaultNamespace().AllNamespaces(d.Namespace == "").
		LabelSelectorParam("").
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

	var obj runtime.Object

	// render `resource.Info` into `runtime.Object`
	{
		// we have zero or multple items, so coerce all items into a list.
		// we don't want an *unstructured.Unstructured list yet, as we
		// may be dealing with non-unstructured objects. Compose all items
		// into an corev1.List, and then decode using an unstructured scheme.
		list := corev1.List{
			TypeMeta: metav1.TypeMeta{
				Kind:       "List",
				APIVersion: "v1",
			},
			ListMeta: metav1.ListMeta{},
		}
		for _, info := range infos {
			list.Items = append(list.Items, runtime.RawExtension{Object: info.Object})
		}

		listData, err := json.Marshal(list)
		if err != nil {
			return nil, err
		}

		converted, err := runtime.Decode(unstructured.UnstructuredJSONScheme, listData)
		if err != nil {
			return nil, err
		}

		obj = converted
	}

	// take the items and create a new list for display
	list := &unstructured.UnstructuredList{
		Object: map[string]interface{}{
			"kind":       "List",
			"apiVersion": "v1",
			"metadata":   map[string]interface{}{},
		},
	}

	// render runtime.Object into `UnstructuredList`
	{
		items, err := meta.ExtractList(obj)
		if err != nil {
			return nil, err
		}

		if listMeta, err := meta.ListAccessor(obj); err == nil {
			list.Object["metadata"] = map[string]interface{}{
				"resourceVersion": listMeta.GetResourceVersion(),
			}
		}

		for _, item := range items {
			list.Items = append(list.Items, *item.(*unstructured.Unstructured))
		}
	}
	return list, nil
}
