package cli

var (
	// HelpFlag is the flag to display the App's help text
	HelpFlag = &Flag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "show help",
	}

	// VersionFlag is the flag to display the App's version text
	VersionFlag = &Flag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "print the version",
	}
)

// Flag is a boolean flag that gets passed down to the action called.
type Flag struct {
	// Name is the name of this flag.
	Name string
	// Aliases is the list of alternate names to enable the flag.
	Aliases []string
	// Usage is a brief description of what the flag is for.
	Usage string
}

// HasName returns true if name matches the flag's name or its aliases
func (f *Flag) HasName(name string) bool {
	aliases := append([]string{f.Name}, f.Aliases...)
	for _, alias := range aliases {
		if name == alias {
			return true
		}
	}

	return false
}
