package common

import "errors"

type KSQLERROR = string

const (
	unsupported KSQLERROR = "unsupported"
)

func Unsupported() error {
	return errors.New(unsupported)
}
