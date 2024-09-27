package cmd

import (
	"fmt"
	"gin.go.dev/pkg/config"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

var (
	cfg     *config.Config
	cfgFile string
	cfgErr  error
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "The main app command",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Go version: %s\n", runtime.Version())
		fmt.Printf("OS: %s\n", runtime.GOOS)
		fmt.Printf("Arch: %s\n", runtime.GOARCH)
		fmt.Printf("CPUs: %d\n\n", runtime.NumCPU())
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is config.toml)")
	rootCmd.AddCommand(cmdServer)
	rootCmd.AddCommand(cmdCreateUser)
	rootCmd.AddCommand(cmdSetPassword)
	rootCmd.AddCommand(cmdMigrate)
}

func initConfig() {
	if cfgFile != "" {
		cfg, cfgErr = config.FromPath(cfgFile)
	} else {
		cfg, cfgErr = config.FromPath("config.toml")
	}

	if cfgErr != nil {
		fmt.Printf("Can't read config: %v\n", cfgErr)
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
