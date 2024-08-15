package cmd

import (
	"fmt"
	"gin.go.dev/internal/config"
	"github.com/spf13/cobra"
	"log"
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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Go version: ", runtime.Version())
		fmt.Println("Operating System: ", runtime.GOOS)
		fmt.Println("CPU Architecture: ", runtime.GOARCH)
		fmt.Println("Number of CPUs: ", runtime.NumCPU())
		fmt.Println("GOROOT: ", runtime.GOROOT())
		fmt.Println("GOPATH: ", os.Getenv("GOPATH"))
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is config.toml)")
	rootCmd.AddCommand(cmdServer)
	rootCmd.AddCommand(cmdCreateUser)
	rootCmd.AddCommand(cmdSetPassword)
}

func initConfig() {
	if cfgFile != "" {
		cfg, cfgErr = config.FromPath(cfgFile)
	} else {
		cfg, cfgErr = config.FromPath("config.toml")
	}

	if cfgErr != nil {
		log.Fatalf("Can't read config: %v\n", cfgErr)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
