package gomake

import "errors"

var (
	// ErrNoSuchTarget is returned if a Gomakefile is ran with an unknown target.
	ErrNoSuchTarget = errors.New("no such target")
)

// Gomakefile is a Makefile representation for gophers.
type Gomakefile struct {
	// Targets is the map of target names to Rules.
	Targets map[string]*Rule
}

// NewGomakefile initializes a Gomakefile that can rebuild itself.
func NewGomakefile() *Gomakefile {
	return &Gomakefile{
		Targets: make(map[string]*Rule),
	}
}

// AddRule creates a new rule and adds it to the Gomakefile.
func (g *Gomakefile) AddRule(target string, dependencies []*Rule, evaluate func() error) *Rule {
	rule := NewRule(target, dependencies, evaluate)
	g.Targets[target] = rule
	return rule
}

// Make makes the target rule and its dependencies.
func (g *Gomakefile) Make(target string) map[string]error {
	rule, ok := g.Targets[target]
	if !ok {
		return map[string]error{
			target: ErrNoSuchTarget,
		}
	}

	return Evaluate(rule)
}
