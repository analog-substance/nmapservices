package cmd

import (
	"github.com/analog-substance/nmapservices"
	"github.com/spf13/cobra"
)

// tcpCmd represents the tcp command
var tcpCmd = &cobra.Command{
	Use:   "tcp",
	Short: "Retrieve TCP nmap service port data",
	Run: func(cmd *cobra.Command, args []string) {
		services, _ := cmd.Flags().GetStringSlice("service")
		topPorts, _ := cmd.Flags().GetInt("top")

		if topPorts > 0 {
			portInfos = nmapservices.TopTCPPorts(topPorts, services...)
		} else {
			portInfos = nmapservices.Ports(nmapservices.TCP, services...)
		}
	},
}

func init() {
	rootCmd.AddCommand(tcpCmd)
}
