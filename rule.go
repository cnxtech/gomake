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
		wg sync.WaitGroup // Waits for all goroutines in rule's dependency graph
		mu sync.Mutex     // Protects errs map
	)
	errs := make(map[*Rule]chan error)

	queue := list.New()
	queue.PushBack(root)
	for elem := queue.Front(); elem != nil; elem = elem.Next() {

		rule := elem.Value.(*Rule)

		// Skip if visited already
		_, ok := errs[rule]
		if ok {
			continue
		}

		// Mark as visited and create result channel
		errs[rule] = make(chan error, 1)

		// Add dependencies to rules to visit
		for _, dependency := range rule.Dependencies {
			queue.PushBack(dependency)
		}

		wg.Add(1)
		go func(rule *Rule) {
			defer wg.Done()
			// Wait for dependencies to be evaluated
			for _, dependency := range rule.Dependencies {
				err := <-errs[dependency]
				mu.Lock()
				errs[dependency] <- err
				mu.Unlock()

				// If any dependency returns err, exit early
				if err != nil {
					mu.Lock()
					errs[rule] <- nil
					mu.Unlock()
					return
				}
			}

			err := rule.Evaluate()
			mu.Lock()
			errs[rule] <- err
			mu.Unlock()
		}(rule)
	}

	wg.Wait()

	// Build results map
	results := make(map[string]error)
	for rule, err := range errs {
		results[rule.Target] = <-err
	}

	return results
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
