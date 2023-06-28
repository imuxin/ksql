package runtime

import (
	"context"
	"fmt"
	"reflect"

	"github.com/reactivex/rxgo/v2"
	"github.com/samber/lo"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/imuxin/ksql/pkg/ext"
	"github.com/imuxin/ksql/pkg/ext/abs"
	"github.com/imuxin/ksql/pkg/ext/kube"
	"github.com/imuxin/ksql/pkg/parser"
)

func convert[T any](list []abs.Object) ([]T, error) {
	_r, err := rxgo.Just(list)().
		Map(func(_ context.Context, i interface{}) (interface{}, error) {
			var t T
			o := reflect.New(reflect.TypeOf(t)).Interface() // o's type is *T
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(i.(abs.Object), o)
			return o, err
		}).
		ToSlice(0)

	if err != nil {
		return nil, err
	}

	result := make([]T, len(_r))
	for i, item := range _r {
		result[i] = *item.(*T)
	}
	return result, nil
}

func (r *RunnableImpl[T]) initPlugin(table, namespace string, k8sFilters KubernetesFilters) error {
	names, selector, err := k8sFilters.Convert()
	if err != nil {
		return err
	}

	r.plugin = ext.NewPlugin(table, namespace, names, r.restConfig, selector)
	return nil
}

func (r *RunnableImpl[T]) initWhereFilter(whereExpr *parser.WhereExpr) {
	r.whereFilter = CompileWhereFilter(whereExpr)
}

func (r *RunnableImpl[T]) list(table, namespace string, k8sFilters []*parser.KubernetesFilter, whereExpr *parser.WhereExpr) ([]abs.Object, error) {
	r.initWhereFilter(whereExpr)
	err := r.initPlugin(table, namespace, k8sFilters)
	if err != nil {
		return nil, err
	}

	list, err := r.plugin.Download()
	if err != nil {
		return nil, err
	}

	return lo.Filter(list, func(item abs.Object, _ int) bool {
		return r.whereFilter.Filter(item)
	}), nil
}

func (r *RunnableImpl[T]) List() ([]T, error) {
	list, err := r.list(
		r.ksql.Select.From.Table,
		r.ksql.Select.Namespace,
		r.ksql.Select.KubernetesFilters,
		r.ksql.Select.Where,
	)
	if err != nil {
		return nil, err
	}

	return convert[T](list)
}

func (r *RunnableImpl[T]) Delete() ([]T, error) {
	list, err := r.list(
		r.ksql.Delete.From.Table,
		r.ksql.Delete.Namespace,
		r.ksql.Delete.KubernetesFilters,
		r.ksql.Delete.Where,
	)
	if err != nil {
		return nil, err
	}
	results, err := r.plugin.Delete(list)
	if err != nil {
		return nil, err
	}
	return convert[T](results)
}

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

	list, err := kube.Describer{
		Tables: tables,
	}.Desc()
	if err != nil {
		return nil, err
	}

	_r, err := rxgo.Just(list)().
		Map(func(_ context.Context, i interface{}) (interface{}, error) {
			var t T
			o := reflect.New(reflect.TypeOf(t)).Interface() // o's type is *T
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(i.(abs.Object), o)
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
