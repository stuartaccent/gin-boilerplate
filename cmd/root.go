package cmd

import (
	"fmt"
	"gin.go.dev/internal/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
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
		headers := []string{"GO VERSION", "OPERATING SYSTEM", "CPU ARCH", "NUM CPUS"}
		row := []string{
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH,
			fmt.Sprintf("%d", runtime.NumCPU()),
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
		fmt.Fprintln(w, strings.Join(headers, "\t"))
		fmt.Fprintln(w, strings.Join(row, "\t"))
		w.Flush()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is config.toml)")
	rootCmd.AddCommand(cmdServer)
	rootCmd.AddCommand(cmdMonitor)
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
