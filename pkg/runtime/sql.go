package runtime

import (
	"context"
	"fmt"
	"reflect"

	"github.com/reactivex/rxgo/v2"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/imuxin/ksql/pkg/ext"
	extkube "github.com/imuxin/ksql/pkg/ext/kube"
	"github.com/imuxin/ksql/pkg/parser"
)

func (r *RunnableImpl[T]) initDownloader(table, namespace string, k8sFilters []*parser.KubernetesFilter) error {
	names := make([]string, 0)
	selector := labels.NewSelector()
	for _, item := range k8sFilters {
		switch {
		case item.Label != nil:
			r, err := item.Label.IntoRequirement()
			if err != nil {
				return err
			}
			selector = selector.Add(*r)
		case item.Name != nil:
			names = append(names, *item.Name)
		}
	}

	r.downloader = extkube.APIServerDownloader{
		RestConfig: r.restConfig,
		Table:      table,
		Namespace:  namespace,
		Names:      names,
		Selector:   selector,
	}
	return nil
}

func (r *RunnableImpl[T]) initWhereFilter(whereExpr *parser.WhereExpr) {
	r.whereFilter = ext.CompileWhereFilter(whereExpr)
}

func (r *RunnableImpl[T]) List() ([]T, error) {
	r.initWhereFilter(r.ksql.Select.Where)
	err := r.initDownloader(r.ksql.Select.From.Table, r.ksql.Select.Namespace, r.ksql.Select.KubernetesFilters)
	if err != nil {
		return nil, err
	}

	list, err := r.downloader.Download()
	if err != nil {
		return nil, err
	}
	_r, err := rxgo.Just(list)().
		Filter(func(i interface{}) bool {
			return r.whereFilter.Filter(i)
		}).
		Map(func(_ context.Context, i interface{}) (interface{}, error) {
			var t T
			o := reflect.New(reflect.TypeOf(t)).Interface() // o's type is *T
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(i.(ext.Object), o)
			return o, err
		}).
		ToSlice(len(list))
	if err != nil {
		return nil, err
	}

	result := make([]T, len(_r))
	for i, item := range _r {
		result[i] = *item.(*T)
	}
	return result, nil
}

func (r *RunnableImpl[T]) Delete() {}

func (r *RunnableImpl[T]) Desc() ([]T, error) {
	sql := fmt.Sprintf("SELECT * FROM crd NAME %s", r.ksql.Desc.Table)
	ksql, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}
	rr := RunnableImpl[apiextensionsv1.CustomResourceDefinition]{
		ksql:       ksql,
		restConfig: r.restConfig,
	}

	tables, err := rr.List()
	if err != nil {
		return nil, err
	}
	if len(tables) == 0 {
		return nil, nil
	}

	list, err := extkube.Describer{
		Tables: tables,
	}.Desc()
	if err != nil {
		return nil, err
	}

	_r, err := rxgo.Just(list)().
		Map(func(_ context.Context, i interface{}) (interface{}, error) {
			var t T
			o := reflect.New(reflect.TypeOf(t)).Interface() // o's type is *T
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(i.(ext.Object), o)
			return o, err
		}).
		ToSlice(len(list))
	if err != nil {
		return nil, err
	}

	result := make([]T, len(_r))
	for i, item := range _r {
		result[i] = *item.(*T)
	}

	return result, nil
}
