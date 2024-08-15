package cmd

import (
	"gin.go.dev/internal/middleware"
	"github.com/spf13/cobra"
	"time"
)

var cmdMonitor = &cobra.Command{
	Use:   "monitor",
	Short: "Run the server and monitor the requests replacing the standard server logging",
	Run: func(cmd *cobra.Command, args []string) {
		go func() {
			for {
				<-time.After(time.Second * 1)
				middleware.MetricsResults.WriteMetrics()
			}
		}()
		cmdServer.Run(cmd, args)
	},
}
