package runtime

type Filter interface {
	Filter(i any) bool
}

var _ Filter = &JsonPathFilter{}

type JsonPathFilter struct{}

func (f JsonPathFilter) Filter(i any) bool {
	return true
}
