package cli

import "testing"

func TestHasName(t *testing.T) {
	flag := &Flag{
		Name:    "help",
		Aliases: []string{"h"},
	}

	if flag.HasName("help ") {
		t.Errorf("Expected flag to return false on unknown name")
	}

	if !flag.HasName("help") {
		t.Errorf("Expected flag to return true on its name")
	}

	if !flag.HasName("h") {
		t.Errorf("Expected flag to return true on its alias")
	}
}

func TestNameForAlias(t *testing.T) {
	flags := Flags{
		{
			Name:    "help",
			Aliases: []string{"h"},
		},
		{
			Name:    "version",
			Aliases: []string{"h"},
		},
	}

	name := flags.NameForAlias("help ")
	if name != "" {
		t.Errorf("Expected no match for unknown alias")
	}

	name = flags.NameForAlias("help")
	if name != "help" {
		t.Errorf("Expected help but got %s", name)
	}

	name = flags.NameForAlias("h")
	if name != "help" {
		t.Errorf("Expected help but got %s", name)
	}
}
