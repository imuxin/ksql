package runtime

import (
	"fmt"
	"strings"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/apiserver/schema"
)

var _ Runnable[any] = &DESCRunnableImpl[any]{}

type DESCRunnableImpl[T any] struct {
	Tables []apiextensionsv1.CustomResourceDefinition
}

type Schema struct {
	Version string `json:"version"`
	Spec    string `json:"spec"`
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

func (r DESCRunnableImpl[T]) Run() ([]T, error) {
	for _, item := range r.Tables[0].Spec.Versions {
		in := item.Schema.OpenAPIV3Schema
		out := &apiextensions.JSONSchemaProps{}
		if err := apiextensionsv1.Convert_v1_JSONSchemaProps_To_apiextensions_JSONSchemaProps(in, out, nil); err != nil {
			return nil, err
		}
		r, _ := schema.NewStructural(out)
		_ = r

		for k, v := range r.Properties {
			TODO(k, v)
		}
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
	return nil, nil
}

func TODO(key string, s schema.Structural) {
	fmt.Println("//", s.Generic.Description)
	switch strings.ToLower(s.Generic.Type) {
	case "object":
		fmt.Println(key, "struct{")
		for k, v := range s.Properties {
			TODO(k, v)
		}
	case "array":
		fmt.Println(key, "[]struct{")
		for k, v := range s.Items.Properties {
			TODO(k, v)
		}
	default:
		fmt.Println(key, s.Type)
	}
}
