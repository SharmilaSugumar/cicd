package yamlparser

import (
	"fmt"

	"go.yaml.in/yaml/v2"
)

type PipelineYAML struct {
	Name string             `yaml:"name"`
	Env  map[string]string  `yaml:"env"`
	Jobs map[string]JobYAML `yaml:"jobs"`
}

type JobYAML struct {
	RunsOn string     `yaml:"runs-on"`
	Steps  []StepYAML `yaml:"steps"`
	Needs  []string   `yaml:"needs"`
}

type StepYAML struct {
	Name    string `yaml:"name"`
	Command string `yaml:"run"`
	Uses    string `yaml:"uses"`
}

// Parse extracts the structured GitHub-Actions-style pipeline instructions from raw YAML.
func Parse(content []byte) (*PipelineYAML, error) {
	var pipeline PipelineYAML
	if err := yaml.Unmarshal(content, &pipeline); err != nil {
		return nil, fmt.Errorf("failed to parse pipeline yaml: %w", err)
	}
	return &pipeline, nil
}
