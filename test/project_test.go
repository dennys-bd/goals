package test

import (
	"testing"

	"github.com/dennys-bd/goals/cmd"
)

type ProjectTest struct {
}

func TestRecriateProject(t *testing.T) {
	p := cmd.RecreateProjectFromGoals("../lib/Goals.toml")
	t.Errorf("Ae: %v", p)
	// println(s)
}
