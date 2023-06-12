package ext

type Object = map[string]interface{}

type Downloader interface {
	Download() ([]Object, error)
}

type Describer interface {
	Desc() ([]Object, error)
}
