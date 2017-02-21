package gomake

import (
	"os"
	"os/exec"
	"testing"
)

func TestGomake(t *testing.T) {
	gomakefile := NewGomakefile()

	gomakefile.AddRule("test", "", nil, func() error {
		build := exec.Command("go", "test", "github.com/goodsport/gomake/pkg/dependency")
		build.Stdout = os.Stdout
		build.Stderr = os.Stderr

		return build.Run()
	})

	err := Gomake(gomakefile).Run([]string{"test"})
	if err != nil {
		t.Errorf("Unexpected err %s", err)
	}
}
