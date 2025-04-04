package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if expected != actual {
		t.Errorf("expected: %v, got: %v", expected, actual)
	}
}

func StringContains(t *testing.T, s, substr string) {
	t.Helper()

	if !strings.Contains(s, substr) {
		t.Errorf("got: %s, expected to contain: %s", s, substr)
	}
}
