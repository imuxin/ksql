package abs

import (
	"github.com/imuxin/ksql/pkg/parser"
	"github.com/imuxin/ksql/pkg/pretty"
)

type Object = map[string]interface{}

type Plugin interface {
	Download() ([]Object, error)
	Delete([]Object) ([]Object, error)
	Columns(*parser.KSQL) []pretty.PrintColumn
}

type Describer interface {
	Desc() ([]Object, error)
}
