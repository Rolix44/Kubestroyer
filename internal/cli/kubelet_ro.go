package cli

import (
	"fmt"

	"github.com/Rolix44/Kubestroyer/internal/display"
	"github.com/Rolix44/Kubestroyer/internal/exploit"
	"github.com/spf13/cobra"
)

func newKubeletROCmd(opts *RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubelet-ro",
		Short: "Check Kubelet read-only port (10255) for unauthenticated endpoint exposure",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, target := range opts.Targets {
				err := exploit.CheckReadOnly(cmd.Context(), exploit.KubeletROOptions{
					Target:      target,
					HTTPTimeout: opts.HTTPTimeout,
				})
				if err != nil {
					fmt.Printf("kubelet-ro error for %s: %v\n", target, err)
				}
				fmt.Println(display.Split)
			}
			return nil
		},
	}
	return cmd
}
