package cli

import "testing"

func TestHasName(t *testing.T) {
	flag := &Flag{
		Name:    "help",
		Aliases: []string{"h"},
	}

	if !flag.HasName("help") {
		t.Errorf("Expected flag to return true on its name")
	}

	if !flag.HasName("h") {
		t.Errorf("Expected flag to return true on its alias")
	}

	if flag.HasName("help ") {
		t.Errorf("Expected flag to return false on unknown name")
	}
}
