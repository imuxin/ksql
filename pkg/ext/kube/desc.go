package kube

import (
	"strings"

	"github.com/fatih/color"
	"github.com/samber/lo"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/apiserver/schema"

	"github.com/imuxin/ksql/pkg/ext/abs"
	"github.com/imuxin/ksql/pkg/util"
)

var _ abs.Describer = &Describer{}

type Describer struct {
	Tables []apiextensionsv1.CustomResourceDefinition
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

func (r Describer) Desc() ([]abs.Object, error) {
	var result []abs.Object
	if len(r.Tables) == 0 {
		return nil, nil
	}
	for _, item := range r.Tables[0].Spec.Versions {
		out := &apiextensions.JSONSchemaProps{}
		if err := apiextensionsv1.Convert_v1_JSONSchemaProps_To_apiextensions_JSONSchemaProps(item.Schema.OpenAPIV3Schema, out, nil); err != nil {
			return nil, err
		}
		s, err := schema.NewStructural(out)
		if err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{
			"version": item.Name,
			"spec":    DeSecrializer(r.Tables[0].Spec.Names.Kind, *s),
		})
	}
	return result, nil
}

func DeSecrializer(kind string, s schema.Structural) string {
	var lines []string
	inner(color.YellowString("type")+" "+color.RedString(kind), s, 0, &lines)
	return strings.Join(lines, "\n")
}

func inner(key string, s schema.Structural, depth int, lines *[]string) {
	var tab = strings.Repeat(strings.Repeat(" ", 4), depth)

	// Description
	for _, item := range util.WrapText(s.Generic.Description, 80) {
		*lines = append(*lines, tab+color.WhiteString("// "+item))
	}

	t := strings.ToLower(s.Generic.Type)
	switch t {
	case "object":
		if len(s.Properties) == 0 {
			t = "interface{}"
			if s.AdditionalProperties != nil && s.AdditionalProperties.Structural.Type != "" {
				t = s.AdditionalProperties.Structural.Type
			}
			*lines = append(*lines, tab+key+" "+color.GreenString("map[string]"+t))
			return
		}
		*lines = append(*lines, tab+key+" "+color.GreenString("struct")+" {")
		depth++
		for item := range util.NewSortRange(s.Properties).Iter() {
			inner(color.MagentaString(item.Key), item.Value, depth, lines)
		}
		*lines = append(*lines, tab+"}")
	case "array":
		if lo.Contains([]string{"object", "array"}, s.Items.Generic.Type) {
			*lines = append(*lines, tab+key+" []"+color.GreenString("struct")+" {")
			depth++
			for item := range util.NewSortRange(s.Items.Properties).Iter() {
				inner(color.MagentaString(item.Key), item.Value, depth, lines)
			}
			*lines = append(*lines, tab+"}")
		} else {
			*lines = append(*lines, tab+key+" []"+color.GreenString(s.Items.Generic.Type))
		}
	case "boolean":
		t = "bool"
		*lines = append(*lines, tab+key+" "+color.GreenString(t))
	case "integer":
		t = "int"
		*lines = append(*lines, tab+key+" "+color.GreenString(t))
	default:
		*lines = append(*lines, tab+key+" "+color.GreenString(t))
	}
}
