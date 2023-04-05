package cmd

import (
	"github.com/analog-substance/nmapservices"
	"github.com/spf13/cobra"
)

// udpCmd represents the udp command
var udpCmd = &cobra.Command{
	Use:   "udp",
	Short: "Retrieve UDP nmap service port data",
	Run: func(cmd *cobra.Command, args []string) {
		services, _ := cmd.Flags().GetStringSlice("service")
		topPorts, _ := cmd.Flags().GetInt("top")

		if topPorts > 0 {
			portInfos = nmapservices.TopUDPPorts(topPorts, services...)
		} else {
			portInfos = nmapservices.Ports(nmapservices.UDP, services...)
		}
	},
}

func init() {
	rootCmd.AddCommand(udpCmd)
}
