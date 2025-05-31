package build

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Project      string              `yaml:"project"`
	Version      string              `yaml:"version"`
	Zip          bool                `yaml:"zip"`
	Build        BuildSettings       `yaml:"build"`
	Dependencies DependencyConfig    `yaml:"dependencies"`
	Tasks        map[string][]string `yaml:"tasks"`
}

type BuildSettings struct {
	Default BuildTarget            `yaml:"default"`
	Targets map[string]BuildTarget `yaml:"targets"`
}

type BuildTarget struct {
	OS        string   `yaml:"os"`
	Arch      string   `yaml:"arch"`
	Output    string   `yaml:"output"`
	Ldflags   string   `yaml:"ldflags,omitempty"`
	Tags      []string `yaml:"tags,omitempty"`
	Cgo       bool     `yaml:"cgo,omitempty"`
	Platforms []string `yaml:"platforms,omitempty"`
}

type DependencyConfig struct {
	Check []string `yaml:"check"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	// Process templates in output paths and ldflags
	if err := processTemplates(config); err != nil {
		return nil, fmt.Errorf("error processing templates: %w", err)
	}

	return config, nil
}

func processTemplates(config *Config) error {
	// Process default target
	if err := processTargetTemplates(&config.Build.Default, config); err != nil {
		return fmt.Errorf("error processing default target templates: %w", err)
	}

	// Process all targets
	for name, target := range config.Build.Targets {
		if err := processTargetTemplates(&target, config); err != nil {
			return fmt.Errorf("error processing target %s templates: %w", name, err)
		}
		config.Build.Targets[name] = target
	}

	// Process task templates
	for taskName, commands := range config.Tasks {
		for i, cmd := range commands {
			processed, err := executeTemplate(cmd, config)
			if err != nil {
				return fmt.Errorf("error processing template in task %s: %w", taskName, err)
			}
			config.Tasks[taskName][i] = processed
		}
	}

	return nil
}

func processTargetTemplates(target *BuildTarget, config *Config) error {
	// Process output template
	if target.Output != "" {
		output, err := executeTemplate(target.Output, config)
		if err != nil {
			return err
		}
		target.Output = output
	}

	// Process ldflags template
	if target.Ldflags != "" {
		ldflags, err := executeTemplate(target.Ldflags, config)
		if err != nil {
			return err
		}
		target.Ldflags = ldflags
	}

	return nil
}

func executeTemplate(templateStr string, config *Config) (string, error) {
	if !strings.Contains(templateStr, "{{") {
		return templateStr, nil
	}

	tmpl, err := template.New("config").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return buf.String(), nil
}
