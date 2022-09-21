package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var scanNode = false

func send_request(target string, port int, endpoint string) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(target + ":" + strconv.Itoa(port) + endpoint)
	if err != nil {
		fmt.Printf("Port %d not open on %s \n", port, target)
	} else {
		fmt.Printf("Port: %d, status: %d \n", port, resp.StatusCode)
		defer resp.Body.Close()
	}
}

func check_ports(target string) {
	endpoint := "/"

	knownPort := []int{443, 2379, 6666, 4194, 6443, 8443, 8080, 10250, 10255, 10256, 9099, 6782, 6783, 6784, 44134}

	if scanNode {
		for i := 30000; i <= 32767; i++ {
			knownPort = append(knownPort, i)
		}
	}

	for _, port := range knownPort {
		if port == 10250 {
			targetHttps := strings.Replace(target, "http", "https", 1)
			endpoint = "/metrics"
			send_request(targetHttps, port, endpoint)
		} else {
			send_request(target, port, endpoint)
		}

	}
}

func main() {
	target := "http://localhost"
	check_ports(target)
}
