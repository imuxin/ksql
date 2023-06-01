package runtime

import (
	"context"
	"reflect"

	"github.com/reactivex/rxgo/v2"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/imuxin/ksql/pkg/ext"
	"github.com/imuxin/ksql/pkg/pretty"
)

var (
	DefaultDatabase = ""
)

var _ Runnable[any] = &RunnableImpl[any]{}

type Runnable[T any] interface {
	Run() ([]T, error)
	RunLikeSQL() ([]pretty.PrintColumn, []T, error)
}

type RunnableImpl[T any] struct {
	Downloader   ext.Downloader
	Filter       Filter
	PrintColumns []pretty.PrintColumn
}

func (r RunnableImpl[T]) Run() ([]T, error) {
	list, err := r.Downloader.Download()
	if err != nil {
		return nil, err
	}
	_r, err := rxgo.Just(list)().
		Filter(func(i interface{}) bool {
			return r.Filter.Filter(i)
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

func (r RunnableImpl[T]) RunLikeSQL() ([]pretty.PrintColumn, []T, error) {
	result, err := r.Run()
	return r.PrintColumns, result, err
}
