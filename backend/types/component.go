package types

import "os/exec"

// Component describes a component in a system
type Component struct {
	Name                     string    `yaml:"name" json:"name"`
	Description              string    `yaml:"description" json:"description"`
	RunCmd                   string    `yaml:"runCmd" json:"runCmd"`
	CurrentlyRunningInstance *exec.Cmd `yaml:"-" json:"currentlyRunningInstance"`
}
