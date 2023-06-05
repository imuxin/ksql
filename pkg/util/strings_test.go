package util

import (
	"reflect"
	"testing"
)

func TestWrapText(t *testing.T) {
	text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce varius enim non tellus elementum, vitae auctor metus facilisis."
	lineLength := 20
	expected := []string{
		"Lorem ipsum dolor",
		"sit amet,",
		"consectetur",
		"adipiscing elit.",
		"Fusce varius enim",
		"non tellus",
		"elementum, vitae",
		"auctor metus",
		"facilisis.",
	}

	result := WrapText(text, lineLength)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result.\nExpected: %v\nGot: %v", expected, result)
	}
}
