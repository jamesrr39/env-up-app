package types

import (
	"errors"
	"fmt"
)

type Environment struct {
	Name       string       `yaml:"name" json:"name"`
	Components []*Component `yaml:"components" json:"components"`
}

func (e *Environment) Validate() error {
	if e.Name == "" {
		return errors.New("the environment needs a name")
	}
	componentNameMap := make(map[string]bool)

	for _, component := range e.Components {
		if component.Name == "" {
			return errors.New("all components need a name")
		}

		isNameAlreadyTaken := componentNameMap[component.Name]
		if isNameAlreadyTaken {
			return fmt.Errorf("the component name %q is already taken by another component in this environment", component.Name)
		}
		componentNameMap[component.Name] = true
	}

	return nil
}
