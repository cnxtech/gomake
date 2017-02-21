package gomake

import (
	"errors"
	"testing"
)

func TestAddRule(t *testing.T) {
	gomakefile := NewGomakefile()
	var expected = errors.New("expected")
	gomakefile.AddRule("target", "", nil, func() error {
		return expected
	})

	rule, ok := gomakefile.Targets["target"]
	if !ok {
		t.Errorf("Expected gomakefile to have target")
	}

	actual := rule.Evaluate()
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestMake(t *testing.T) {
	gomakefile := NewGomakefile()
	results := gomakefile.Make("target")
	if results["target"] == nil {
		t.Errorf("Unknown target doesn't return error")
	}

	var expected = errors.New("expected")
	rule := gomakefile.AddRule("target", "", nil, func() error {
		return expected
	})

	actual := rule.Evaluate()
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
