package cli

import "testing"

func TestRun(t *testing.T) {
	app := &App{}
	err := app.Run([]string{"gomake"})
	if err != nil {
		t.Errorf("Unexpected err")
	}
}

func TestShowHelp(t *testing.T) {
	app := &App{}
	err := app.ShowHelp()
	if err != nil {
		t.Errorf("Unexpected err")
	}
}
