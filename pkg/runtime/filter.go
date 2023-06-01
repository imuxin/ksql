package runtime

type Filter interface {
	Filter(i any) bool
}
