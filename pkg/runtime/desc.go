package runtime

import (
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/apiserver/schema"

	"github.com/alecthomas/repr"
	"github.com/imuxin/ksql/pkg/pretty"
)

var _ Runnable[any] = &DESCRunnableImpl[any]{}

type DESCRunnableImpl[T any] struct {
	Tables []apiextensionsv1.CustomResourceDefinition
}

func (r DESCRunnableImpl[T]) Run() ([]T, error) {
	return nil, nil
}

/*
will output like this:

	struct{
		spec struct{
			// description here
			name string // required

			// description here
			age int

			array []string

			xxx struct{
				hah bool // required

				aa []string // required
			}
		}
	}
*/
func (r DESCRunnableImpl[T]) RunLikeSQL() ([]pretty.PrintColumn, []T, error) {
	for _, item := range r.Tables[0].Spec.Versions {
		in := item.Schema.OpenAPIV3Schema
		out := &apiextensions.JSONSchemaProps{}
		if err := apiextensionsv1.Convert_v1_JSONSchemaProps_To_apiextensions_JSONSchemaProps(in, out, nil); err != nil {
			return nil, nil, err
		}
		r, _ := schema.NewStructural(out)
		_ = r
		repr.Println(r)
	}
	// root := r.Tables[0].Spec.Versions[0].Schema.OpenAPIV3Schema
	// repr.Println(schema.NewStructural(root))
	// required := root.Required
	// types := root.Type
	// desc := root.Description
	// if types == "object" {
	// 	for _, item := range root.Properties {
	// 		// 递归 item
	// 	}
	// }
	return nil, nil, nil
}
