package runtime

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/samber/lo"
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

		fmt.Println(DeSecrializer(*r))
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

func DeSecrializer(s schema.Structural) string {
	var lines []string
	inner("", s, -1, &lines)
	return strings.Join(lines, "")
}

func inner(key string, s schema.Structural, depth int, lines *[]string) {
	key = color.MagentaString(key)
	var tab string
	if depth > 0 {
		tab = strings.Repeat("\t", depth)
	}

	*lines = append(*lines, fmt.Sprintln(tab, color.WhiteString("// "+s.Generic.Description)))
	switch strings.ToLower(s.Generic.Type) {
	case "object":
		*lines = append(*lines, fmt.Sprintln(tab, key, "struct{"))
		depth++
		for k, v := range s.Properties {
			inner(k, v, depth, lines)
		}
	case "array":

		if lo.Contains([]string{"object", "array"}, s.Items.Generic.Type) {
			*lines = append(*lines, fmt.Sprintln(tab, key, "[]struct{"))
			depth++
			for k, v := range s.Items.Properties {
				inner(k, v, depth, lines)
			}
		} else {
			*lines = append(*lines, fmt.Sprintln(tab, key, "[]"+color.GreenString(s.Items.Generic.Type)))
			return
		}
	default:
		*lines = append(*lines, fmt.Sprintln(tab, key, color.GreenString(s.Generic.Type)))
	}
}

// func TODO(key string, s schema.Structural, depth int) {
// 	tab := strings.Repeat("\t", depth)
// 	fmt.Println(tab, "//", s.Generic.Description)
// 	switch strings.ToLower(s.Generic.Type) {
// 	case "object":
// 		fmt.Println(tab, key, "struct{")
// 		depth += 1
// 		for k, v := range s.Properties {
// 			TODO(k, v, depth)
// 		}
// 	case "array":
// 		fmt.Println(tab, key, "[]struct{")
// 		depth += 1
// 		for k, v := range s.Items.Properties {
// 			TODO(k, v, depth)
// 		}
// 	default:
// 		fmt.Println(tab, key, s.Type)
// 	}
// }
