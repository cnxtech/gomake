package gomake

import (
	"container/list"
	"fmt"
	"sync"
)

// Rule is a node in a dependency graph.
type Rule struct {
	// Target is the identifier for the rule in the results from Evaluate.
	Target string
	// Description is an optional field describing the rule.
	Description string
	// Dependencies is a list of rules that must be evaluated before this.
	Dependencies []*Rule
	// Evaluate is the arbitrary function to evaluate the rule.
	Evaluate func() error
}

// NewRule initializes a new named Rule with its direct dependencies and
// evaluate function.
func NewRule(target, description string, dependencies []*Rule, evaluate func() error) *Rule {
	return &Rule{
		Target:       target,
		Description:  description,
		Dependencies: dependencies,
		Evaluate:     evaluate,
	}
}

// Evaluate traverses root rule's dependency graph and creates goroutines for
// all rules it visit. Each goroutine will wait for its dependencies to be
// evaluated before evaluating itself, but if any dependency evaluates with an
// error, it will exit early.
func Evaluate(root *Rule) map[string]error {
	var (
		// Waits for all goroutines in rule's dependency graph to finish evaluating
		wg sync.WaitGroup
		// Protects errs map
		mu sync.Mutex
	)

	resultChs := make(map[*Rule]chan error)

	// Stall rule evaluation until all rules have been visited
	mu.Lock()

	queue := list.New()
	queue.PushBack(root)
	for elem := queue.Front(); elem != nil; elem = elem.Next() {
		rule := elem.Value.(*Rule)

		// Skip if visited already
		_, ok := resultChs[rule]
		if ok {
			continue
		}

		// Mark as visited and create result channel
		resultChs[rule] = make(chan error, 1)

		// Add dependencies to rules to visit
		for _, dependency := range rule.Dependencies {
			queue.PushBack(dependency)
		}

		wg.Add(1)
		go func(rule *Rule) {
			defer wg.Done()
			evaluateRule(rule, &mu, resultChs)
		}(rule)
	}

	// Rules can begin evaluating
	mu.Unlock()
	wg.Wait()

	// Build results map
	results := make(map[string]error)
	for rule, err := range resultChs {
		results[rule.Target] = <-err
	}

	return results
}

func evaluateRule(rule *Rule, mu *sync.Mutex, resultChs map[*Rule]chan error) {
	mu.Lock()
	ruleCh := resultChs[rule]
	mu.Unlock()

	// Wait for dependencies to be evaluated
	for _, dependency := range rule.Dependencies {
		mu.Lock()
		dependencyCh := resultChs[dependency]
		mu.Unlock()

		// Grab a copy and return it to the channel so that all its dependents
		// can take a look at its result
		err := <-dependencyCh
		dependencyCh <- err

		// If any dependency returns err, exit early
		if err != nil {
			ruleCh <- nil
			return
		}
	}

	err := rule.Evaluate()
	ruleCh <- err
}

// HandleResults displays all the target errs and returns a combined error.
func HandleResults(results map[string]error) error {
	var errs []error
	for target, err := range results {
		if err != nil {
			fmt.Printf("%s: %s", target, err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", errs)
	}

	return nil
}
