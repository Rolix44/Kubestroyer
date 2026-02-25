package scanner

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

// Config holds scanner configuration.
type Config struct {
	Target      string
	PortMap     map[int]string
	Workers     int
	TCPTimeout  time.Duration
	HTTPTimeout time.Duration
}

// Result holds the result of a successful port probe.
type Result struct {
	Port        int
	Description string
}

// Scan concurrently probes all ports in cfg.PortMap and returns open ones sorted by port number.
func Scan(ctx context.Context, cfg Config) ([]Result, error) {
	if cfg.Workers <= 0 {
		cfg.Workers = 100
	}
	if cfg.TCPTimeout == 0 {
		cfg.TCPTimeout = 5 * time.Second
	}
	if cfg.HTTPTimeout == 0 {
		cfg.HTTPTimeout = 3 * time.Second
	}

	sem := make(chan struct{}, cfg.Workers)
	results := make(chan Result, len(cfg.PortMap))
	var wg sync.WaitGroup

	for port := range cfg.PortMap {
		wg.Add(1)
		sem <- struct{}{}
		go func(p int) {
			defer wg.Done()
			defer func() { <-sem }()
			if probePort(ctx, cfg, p) {
				results <- Result{Port: p, Description: cfg.PortMap[p]}
			}
		}(port)
	}

	go func() { wg.Wait(); close(results) }()

	var out []Result
	for r := range results {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Port < out[j].Port })
	return out, nil
}

func newHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout,
	}
}

// probePort checks whether a port is open.
// Ports 443 and 10250 are probed via HTTPS; all others via TCP dial.
func probePort(ctx context.Context, cfg Config, port int) bool {
	switch port {
	case 443, 10250:
		client := newHTTPClient(cfg.HTTPTimeout)
		endpoint := "/"
		if port == 10250 {
			endpoint = "/metrics"
		}
		url := "https://" + cfg.Target + ":" + strconv.Itoa(port) + endpoint
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return false
		}
		resp, err := client.Do(req)
		if err != nil {
			return false
		}
		resp.Body.Close()
		return true
	default:
		dialer := &net.Dialer{Timeout: cfg.TCPTimeout}
		conn, err := dialer.DialContext(ctx, "tcp", cfg.Target+":"+strconv.Itoa(port))
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}
}
