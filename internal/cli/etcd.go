package cli

import (
	"fmt"

	"github.com/Rolix44/Kubestroyer/internal/display"
	"github.com/Rolix44/Kubestroyer/internal/exploit"
	"github.com/spf13/cobra"
)

func newEtcdCmd(opts *RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "etcd",
		Short: "Read etcd objects anonymously on target(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, target := range opts.Targets {
				err := exploit.ReadObjects(cmd.Context(), exploit.EtcdOptions{
					Target: target,
				})
				if err != nil {
					fmt.Printf("Etcd error for %s: %v\n", target, err)
				}
				fmt.Println(display.Split)
			}
			return nil
		},
	}
	return cmd
}
