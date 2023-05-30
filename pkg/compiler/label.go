package compiler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"

	"github.com/imuxin/ksql/pkg/parser"
)

type LabelCompare parser.Compare

func (c LabelCompare) IntoRequirement() (*labels.Requirement, error) {
	op, err := c.Op()
	if err != nil {
		return nil, err
	}

	var vals []string
	switch op {
	case selection.Exists, selection.DoesNotExist:
	default:
		vals = Value(c.Operation.RHS).Into()
	}
	return labels.NewRequirement(c.LHS, op, vals)
}

func (c LabelCompare) Op() (selection.Operator, error) {
	if c.NOT {
		return "", errors.New("unexpected operator `NOT` before label compare expr")
	}
	switch strings.ToLower(c.Operation.Exists) {
	case "exists":
		return selection.Exists, nil
	case "notexists":
		return selection.DoesNotExist, nil
	}
	switch strings.ToLower(c.Operation.Op) {
	case "=":
		return selection.Equals, nil
	case "==":
		return selection.DoubleEquals, nil
	case "<>", "!=":
		return selection.NotEquals, nil
	case "in":
		return selection.In, nil
	case "notin":
		return selection.NotIn, nil
	case ">=":
		return selection.GreaterThan, nil
	case "<=":
		return selection.LessThan, nil
	}
	return "", fmt.Errorf("unexpected operator `%s` in label compare expr", c.Operation.Op)
}

type Value parser.Value

func (v Value) Into() []string {
	vals := make([]string, 0)
	switch {
	case v.Array != nil:
		vals = v.IntoArray()
	default:
		vals = append(vals, v.IntoSingle())
	}
	return vals
}

func (v Value) IntoSingle() string {
	var val string
	switch {
	case v.Boolean != nil:
		val = strconv.FormatBool(bool(*v.Boolean))
	case v.Null:
		val = ""
	case v.Number != nil:
		val = strconv.FormatFloat(*v.Number, 'f', -1, 64)
	case v.String != nil:
		val = *v.String
	default:
		val = ""
	}
	return val
}

func (v Value) IntoArray() []string {
	vals := make([]string, 0)
	for _, item := range v.Array.Value {
		val := (Value)(*item).IntoSingle()
		vals = append(vals, val)
	}
	return vals
}
