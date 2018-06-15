package types

type Environment struct {
	Name       string       `yaml:"name" json:"name"`
	Components []*Component `yaml:"components" json:"components"`
}
