package types

import "os"

// Component describes a component in a system
type Component struct {
	Name                     string      `yaml:"name" json:"name"`
	Description              string      `yaml:"description" json:"description"`
	RunCmd                   string      `yaml:"runCmd" json:"runCmd"`
	CurrentlyRunningInstance *os.Process `yaml:"-" json:"currentlyRunningInstance"`
}
