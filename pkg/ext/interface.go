package ext

type Object = map[string]interface{}

type Plugin interface {
	Download() ([]Object, error)
	Delete([]Object) error
}

type Describer interface {
	Desc() ([]Object, error)
}
