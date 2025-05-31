package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Executor struct {
	config *Config
}

func NewExecutor(config *Config) *Executor {
	return &Executor{config: config}
}

func (e *Executor) BuildTarget(targetName string) error {

	if targetName == "all" {
		for _, platform := range e.config.Build.Targets[targetName].Platforms {
			err := e.BuildSingleTarget(platform)
			if err != nil {
				return err
			}
		}
	} else {
		return e.BuildSingleTarget(targetName)
	}
	return nil
}

func (e *Executor) BuildSingleTarget(targetName string) error {
	var target BuildTarget
	var ok bool

	if targetName == "default" {
		target = e.config.Build.Default
	} else {
		target, ok = e.config.Build.Targets[targetName]
		if !ok {
			return fmt.Errorf("unknown build target: %s", targetName)
		}
	}

	if err := os.MkdirAll(filepath.Dir(target.Output), 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	cmd := exec.Command("go", "build")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", target.OS),
		fmt.Sprintf("GOARCH=%s", target.Arch),
	)
	if target.Cgo {
		cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
	} else {
		cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
	}

	if target.Ldflags != "" {
		cmd.Args = append(cmd.Args, "-ldflags", target.Ldflags)
	}

	if len(target.Tags) > 0 {
		cmd.Args = append(cmd.Args, "-tags", strings.Join(target.Tags, ","))
	}

	cmd.Args = append(cmd.Args, "-o", target.Output, ".")

	fmt.Printf("Building for %s/%s to %s...\n", target.OS, target.Arch, target.Output)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	if e.config.Zip {
		// rename the output file to include the version
		if err := os.Rename(target.Output, fmt.Sprintf("%s/%s", filepath.Dir(target.Output), e.config.Project)); err != nil {
			return fmt.Errorf("error renaming output file: %w", err)
		}

		cmd = exec.Command("zip", "-j", fmt.Sprintf("%s.zip", target.Output), fmt.Sprintf("%s/%s", filepath.Dir(target.Output), e.config.Project))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error creating zip file: %w", err)
		}
		fmt.Printf("Output zipped to: %s\n", target.Output)
		os.Remove(fmt.Sprintf("%s/%s", filepath.Dir(target.Output), e.config.Project))
	}
	fmt.Printf("Build successful: %s\n", target.Output)
	return nil
}

func (e *Executor) RunTasks(taskName string) error {
	tasks, ok := e.config.Tasks[taskName]
	if !ok {
		return fmt.Errorf("unknown task: %s", taskName)
	}

	for _, task := range tasks {
		fmt.Printf("Running task: %s\n", task)
		cmd := exec.Command("sh", "-c", task)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("task failed: %w", err)
		}
	}
	return nil
}

func (e *Executor) RunPreBuildTasks() error {
	if _, ok := e.config.Tasks["pre-build"]; ok {
		return e.RunTasks("pre-build")
	}
	return nil
}

func (e *Executor) RunPostBuildTasks() error {
	if _, ok := e.config.Tasks["post-build"]; ok {
		return e.RunTasks("post-build")
	}
	return nil
}
