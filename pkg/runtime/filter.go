package runtime

type Filter interface {
	Filter(i any) bool
}

var _ Filter = &JSONPathFilter{}

type JSONPathFilter struct{}

func (f JSONPathFilter) Filter(_ any) bool {
	return true
}
