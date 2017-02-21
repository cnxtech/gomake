package cli

// Command is a subcommand for an App.
type Command struct {
	// Name is the name of the subcommand.
	Name string
	// Description is a brief text about the subcommand.
	Description string
	// Action is the function to call when the command is invoked.
	Action Action
}

// Commands is a sortable list of commands.
type Commands []*Command

// Len returns the length of commands.
func (c Commands) Len() int {
	return len(c)
}

// Len returns whether the command at index i is less than at index j.
func (c Commands) Less(i, j int) bool {
	return c[i].Name < c[j].Name
}

// Swap swaps the commands at index i and j.
func (c Commands) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Commands) ActionForName(name string) Action {
	for _, command := range c {
		if name == command.Name {
			return command.Action
		}
	}

	return nil
}
