package compiler

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"

	extkube "github.com/imuxin/ksql/pkg/ext/kube"
	"github.com/imuxin/ksql/pkg/parser"
	"github.com/imuxin/ksql/pkg/runtime"
)

func Compile[T any](ksql *parser.KSQL, restConfig *rest.Config) (runtime.Runnable[T], error) {
	if ksql.Select.From.Table == "" {
		return nil, nil
	}

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

	d := extkube.APIServerDownloader{
		RestConfig: restConfig,
		Table:      ksql.Select.From.Table,
		Namespace:  ksql.Select.Namespace,
		Names:      names,
		Selector:   selector,
	}
	return runtime.RunnableImpl[T]{
		Downloader: d,
		Filter:     runtime.JSONPathFilter{}, // TODO
	}, nil
}
