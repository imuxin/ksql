package compiler

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/parser"
	"github.com/imuxin/ksql/pkg/runtime"
)

func Compile[T any](ksql *parser.KSQL, restConfig *rest.Config) (runtime.Runnable[T], error) {
	names := make([]string, 0)
	selector := labels.NewSelector()
	for _, item := range ksql.Select.KubernetesFilters {
		switch {
		case item.Label != nil:
			r, err := (LabelCompare)(*item.Label).IntoRequirement()
			if err != nil {
				return nil, err
			}
			selector = selector.Add(*r)
		case item.Name != nil:
			names = append(names, *item.Name)
		}
	}

	d := runtime.APIServerDownloader{
		RestConfig: restConfig,
		Table:      ksql.Select.From.Table,
		Namespace:  ksql.Select.Namespace,
		Names:      names,
		Selector:   selector,
	}
	return runtime.KubernetesRunnable[T]{
		Downloader: d,
		Filter:     runtime.JSONPathFilter{}, // TODO
	}, nil
}
