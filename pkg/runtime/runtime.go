package runtime

import (
	"context"
	"reflect"

	"github.com/reactivex/rxgo/v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	DefaultDatabase = ""
)

var _ Runnable[any] = &KubernetesRunnable[any]{}

type Runnable[T any] interface {
	Run() ([]T, error)
}

type KubernetesRunnable[T any] struct {
	Downloader Downloader
	Filter     Filter
}

func (r KubernetesRunnable[T]) Run() ([]T, error) {
	list, err := r.Downloader.Download()
	if err != nil {
		return nil, err
	}
	_r, err := rxgo.Just(list.Items)().
		Filter(func(i interface{}) bool {
			return r.Filter.Filter(i)
		}).
		Map(func(_ context.Context, i interface{}) (interface{}, error) {
			var t T
			o := reflect.New(reflect.TypeOf(t)).Interface() // o's type is *T
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(i.(unstructured.Unstructured).Object, o)
			return o, err
		}).
		ToSlice(len(list.Items))
	if err != nil {
		return nil, err
	}

	result := make([]T, len(_r))
	for i, item := range _r {
		result[i] = *item.(*T)
	}
	return result, nil
}
