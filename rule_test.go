package gomake

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"testing"
)

func TestEvaluate(t *testing.T) {
	var (
		actual []byte
		// Protects actual byte array
		mu sync.Mutex
	)
	rule1 := NewRule("", "", nil, func() error {
		mu.Lock()
		defer mu.Unlock()
		actual = append(actual, '1')
		return nil
	})
	rule2 := NewRule("", "", nil, func() error {
		mu.Lock()
		defer mu.Unlock()
		actual = append(actual, '1')
		return nil
	})
	rule3 := NewRule("", "", []*Rule{rule1, rule2}, func() error {
		mu.Lock()
		defer mu.Unlock()
		actual = append(actual, '2')
		return nil
	})
	rule4 := NewRule("", "", []*Rule{rule2, rule3}, func() error {
		mu.Lock()
		defer mu.Unlock()
		actual = append(actual, '3')
		return nil
	})

	err := HandleResults(Evaluate(rule4))
	if err != nil {
		t.Errorf("Failed to evaluate", err)
	}

	expected := []byte{'1', '1', '2', '3'}
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected order %s but got %s", expected, actual)
	}
}

func TestEvaluateErr(t *testing.T) {
	var (
		actual []byte
		// Protects actual byte array
		mu sync.Mutex
	)
	rule1 := NewRule("", "", nil, func() error {
		mu.Lock()
		defer mu.Unlock()
		actual = append(actual, '1')
		return nil
	})
	rule2 := NewRule("", "", []*Rule{rule1}, func() error {
		mu.Lock()
		defer mu.Unlock()
		actual = append(actual, '2')
		return nil
	})
	intentional := errors.New("intentional")
	rule3 := NewRule("error", "", []*Rule{rule1}, func() error {
		return intentional
	})
	rule4 := NewRule("", "", []*Rule{rule2, rule3}, func() error {
		mu.Lock()
		defer mu.Unlock()
		actual = append(actual, '3')
		return nil
	})
	rule5 := NewRule("", "", []*Rule{rule3}, func() error {
		mu.Lock()
		defer mu.Unlock()
		actual = append(actual, '4')
		return nil
	})
	rule6 := NewRule("", "", []*Rule{rule4, rule5}, func() error {
		return nil
	})

	results := Evaluate(rule6)
	err, ok := results["error"]
	if !ok {
		t.Errorf("No result for target error")
	}

	if err != intentional {
		t.Errorf("Expected %s but got %s", intentional, err)
	}

	expected := []byte{'1', '2'}
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected order %s but got %s", expected, actual)
	}
}

func TestHandleResults(t *testing.T) {
	results := map[string]error{
		"target": nil,
	}

	err := HandleResults(results)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	expected := errors.New("expected")
	results = map[string]error{
		"target1": nil,
		"target2": expected,
	}

	err = HandleResults(results)
	if err == nil {
		t.Errorf("Expected err")
	}

	if !strings.Contains(err.Error(), expected.Error()) {
		t.Errorf("Expected %s in error message but got %s", expected, err)
	}
}
