package executor

import (
	"testing"

	"github.com/imuxin/ksql/pkg/common"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestUnSupport(t *testing.T) {
	for _, sql := range []string{
		"use xxx",
		"delete",
		"update",
	} {
		_, err := Execute[unstructured.Unstructured](sql, nil)
		assert.EqualError(t, err, common.Unsupported().Error())
	}
}

func TestDESC(_ *testing.T) {
	for _, sql := range []string{
		"desc workflows.core.oam.dev",
	} {
		_, _, _ = ExecuteLikeSQL[unstructured.Unstructured](sql, nil)
	}
}
