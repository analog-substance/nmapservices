package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/analog-substance/nmapservices"
	"github.com/spf13/cobra"
)

var portInfos []nmapservices.PortInfo

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nmapservices",
	Short: "A tool to retrieve nmap service port data",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if portInfos == nil {
			return
		}

		asRanges, _ := cmd.Flags().GetBool("ranges")
		portsOnly, _ := cmd.Flags().GetBool("quite")
		delimiter, _ := cmd.Flags().GetString("delimiter")

		var items []string
		var ports []int
		for _, p := range portInfos {
			ports = append(ports, p.Port)

			if portsOnly {
				items = append(items, fmt.Sprintf("%d", p.Port))
			} else {
				items = append(items, fmt.Sprintf("%d (%s)", p.Port, p.Service))
			}
		}

		if asRanges {
			items = toRanges(ports)
		}

		fmt.Println(strings.Join(items, delimiter))
	},
}

func toRanges(ports []int) []string {
	if len(ports) == 0 {
		return nil
	}

	sort.Ints(ports)

	var ranges []string
	start := 0
	previous := 0

	for _, port := range ports {
		if previous == 0 {
			previous = port
			start = port
			continue
		}

		if port == previous+1 {
			previous = port
			continue
		}

		if start == previous {
			ranges = append(ranges, fmt.Sprintf("%d", start))
		} else {
			ranges = append(ranges, fmt.Sprintf("%d-%d", start, previous))
		}

		previous = port
		start = port
	}

	return ranges
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("quite", "q", false, "Only show port numbers")
	rootCmd.PersistentFlags().Bool("ranges", false, "Output as port ranges. Implies --ports")
	rootCmd.PersistentFlags().StringP("delimiter", "d", "\n", "Delimit ports by specified string")
	rootCmd.PersistentFlags().StringSliceP("service", "s", []string{}, "Output ports for services. Matches if the port's service contains any of the specified service strings")
	rootCmd.PersistentFlags().IntP("top", "t", 0, "Output the top number ports")
}
