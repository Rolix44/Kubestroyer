package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Rolix44/Kubestroyer/internal/display"
	"github.com/Rolix44/Kubestroyer/internal/ports"
	"github.com/Rolix44/Kubestroyer/internal/scanner"
	"github.com/spf13/cobra"
)

func newScanCmd(opts *RootOptions) *cobra.Command {
	var nodePorts bool

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan Kubernetes ports on target(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			portMap := ports.KnownPorts()
			if nodePorts {
				portMap = ports.WithNodePorts()
			}
			for _, target := range opts.Targets {
				if err := runScan(cmd.Context(), target, portMap, opts); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&nodePorts, "node-ports", false, "Scan NodePort range 30000-32767")
	return cmd
}

// runScan executes a port scan for a single target and prints the results.
// Shared by the scan and all subcommands.
func runScan(ctx context.Context, target string, portMap map[int]string, opts *RootOptions) error {
	fmt.Printf("Starting port scan for '%s'...\n\n", target)
	results, err := scanner.Scan(ctx, scanner.Config{
		Target:      target,
		PortMap:     portMap,
		Workers:     opts.Workers,
		TCPTimeout:  opts.TCPTimeout,
		HTTPTimeout: opts.HTTPTimeout,
	})
	if err != nil {
		return err
	}
	if len(results) == 0 {
		fmt.Println("\x1b[1;31mNo open ports found!\x1b[0m")
	} else {
		for _, r := range results {
			fmt.Println("\x1b[1;32m[+]\x1b[0m port " + strconv.Itoa(r.Port) + " open (" + r.Description + ")")
		}
	}
	fmt.Println(display.Split)
	return nil
}
