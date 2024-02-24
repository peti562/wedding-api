package helpers_test

import (
	"testing"

	"github.com/peti562/wedding/pkg/helpers"
)

func TestIsArrayContains(t *testing.T) {
	// test case 1
	arr := []string{"hello", "world", "golang"}
	str := "golang"
	expected := true
	result := helpers.IsArrayContains(arr, str)
	if result != expected {
		t.Errorf("Expected %t but got %t", expected, result)
	}

	// test case 2
	arr = []string{"hello", "world", "golang"}
	str = "java"
	expected = false
	result = helpers.IsArrayContains(arr, str)
	if result != expected {
		t.Errorf("Expected %t but got %t", expected, result)
	}
}
