package pkg

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/Rolix44/Kubestroyer/utils"
)

var openPort []int

func nodeport_scan(target string, port int) {
	servAddr := target + ":" + strconv.Itoa(port)
	_, err := net.Dial("tcp", servAddr)
	if !strings.Contains(err.Error(), "connection refused") {
		openPort = append(openPort, port)
	}

}

func send_http_request(target string, port int, endpoint string) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(target + ":" + strconv.Itoa(port) + endpoint)
	if err == nil {
		openPort = append(openPort, port)
		defer resp.Body.Close()
	}
}

func check_ports(target string) {
	fmt.Printf("Starting port scan for '%s'... \n\n", target)
	openPort = nil

	if utils.ScanNode {
		for port := 30000; port <= 32767; port++ {
			nodeport_scan(target, port)
		}
	}
	target = "http://" + target
	endpoint := "/"

	for port := range utils.KnownPorts {
		if port == 10250 {
			target := strings.Replace(target, "http", "https", 1)
			endpoint = "/metrics"
			send_http_request(target, port, endpoint)
		} else {

			send_http_request(target, port, endpoint)
		}

	}

	if len(openPort) != 0 {
		for _, port := range openPort {
			fmt.Println("\x1b[1;32m[+]\x1b[0m port " + strconv.Itoa(port) + " open (" + utils.KnownPorts[port] + ")")
		}
	} else {
		fmt.Println("\x1b[1;31mNo open ports found !\x1b[0m")
	}
	fmt.Println(utils.Split)

}
