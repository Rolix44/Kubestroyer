package pkg

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Rolix44/Kubestroyer/utils"
)

var openPort []int

func nodeportScan(target string, port int) {
	d := net.Dialer{Timeout: 5}
	servAddr := target + ":" + strconv.Itoa(port)
	_, err := d.Dial("tcp", servAddr)
	if err != nil {
		return
	}
	openPort = append(openPort, port)

}

func sendHttpRequest(target string, port int, endpoint string) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get(target + ":" + strconv.Itoa(port) + endpoint)
	if err == nil {
		openPort = append(openPort, port)
		defer resp.Body.Close()
	}
}

func checkPorts(target string) {
	fmt.Printf("Starting port scan for '%s'... \n\n", target)
	openPort = nil

	if utils.ScanNode {
		for port := 30000; port <= 32767; port++ {
			nodeportScan(target, port)
		}
	}
	target = "http://" + target
	endpoint := "/"

	for port := range utils.KnownPorts {
		if port == 10250 || port == 443 {
			target := strings.Replace(target, "http", "https", 1)
			if port == 10250 {
				endpoint = "/metrics"
			}
			sendHttpRequest(target, port, endpoint)
			endpoint = "/"
		} else {
			sendHttpRequest(target, port, endpoint)
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
