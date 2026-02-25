package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Rolix44/Kubestroyer/internal/display"
	"github.com/spf13/cobra"
)

// RootOptions holds flags shared across all subcommands.
type RootOptions struct {
	RawTarget   string
	Targets     []string
	Workers     int
	TCPTimeout  time.Duration
	HTTPTimeout time.Duration
}

// Execute is the entry point called by main.
func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	opts := &RootOptions{}

	root := &cobra.Command{
		Use:   "kubestroyer",
		Short: "Kubernetes penetration testing tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			display.PrintBanner()
			targets, err := parseTargets(opts.RawTarget)
			if err != nil {
				return err
			}
			opts.Targets = targets
			return nil
		},
	}

	root.PersistentFlags().StringVarP(&opts.RawTarget, "target", "t", "", "Target IP, domain, comma-separated list, or file path")
	root.PersistentFlags().IntVar(&opts.Workers, "workers", 100, "Number of concurrent workers")
	root.PersistentFlags().DurationVar(&opts.TCPTimeout, "tcp-timeout", 5*time.Second, "TCP connection timeout")
	root.PersistentFlags().DurationVar(&opts.HTTPTimeout, "http-timeout", 3*time.Second, "HTTP request timeout")
	root.MarkPersistentFlagRequired("target")

	root.AddCommand(
		newAllCmd(opts),
		newScanCmd(opts),
		newKubeletROCmd(opts),
		newRCECmd(opts),
		newEtcdCmd(opts),
	)

	return root
}

// parseTargets resolves a raw target string into a list of targets.
// If raw points to a file, targets are read line by line; otherwise split on commas.
func parseTargets(raw string) ([]string, error) {
	_, err := os.Stat(raw)
	if err != nil {
		return strings.Split(raw, ","), nil
	}

	file, err := os.Open(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to open target file: %w", err)
	}
	defer file.Close()

	var targets []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		if line := sc.Text(); line != "" {
			targets = append(targets, line)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("error reading target file: %w", err)
	}
	return targets, nil
}
