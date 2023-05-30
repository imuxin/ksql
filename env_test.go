package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

// func currentDir() string {
// 	_, currentFile, _, _ := runtime.Caller(0) //nolint:dogsled
// 	return filepath.Dir(currentFile)
// }

func TestMain(t *testing.T) {
	env := &envtest.Environment{}
	_, err := env.Start()
	assert.NoError(t, err)
}
