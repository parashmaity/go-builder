package build

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cfgFile string
var targetFlag string

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project",
	Long: `Build the project for specific target or the default target if none is specified.
		Available targets are defined in build-config.yml.`,
	Run: func(cmd *cobra.Command, args []string) {
		target := "default"
		if targetFlag != "" {
			target = targetFlag
		}

		config, err := LoadConfig(cfgFile)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		executor := NewExecutor(config)
		if err := executor.RunPreBuildTasks(); err != nil {
			fmt.Printf("Pre-build tasks failed: %v\n", err)
			return
		}

		if err := executor.BuildTarget(target); err != nil {
			fmt.Printf("Build failed: %v\n", err)
			return
		}

		if err := executor.RunPostBuildTasks(); err != nil {
			fmt.Printf("Post-build tasks failed: %v\n", err)
			return
		}
	},
}

func init() {
	BuildCmd.PersistentFlags().StringVar(&cfgFile, "config", "build-config.yaml", "config file (default is build-config.yaml)")
	BuildCmd.Flags().StringVar(&targetFlag, "target", "", "build target to use (default is 'default')")
}
