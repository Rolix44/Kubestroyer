package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func check_ports(domain string) {
	known_port := []int{443, 2379, 6666, 4194, 6443, 8443, 8080, 10250, 10255, 10256, 9099, 44134}

	for _, port := range known_port {
		resp, err := http.Get(domain + ":" + strconv.Itoa(port))
		if err != nil {
			fmt.Printf("Port %d not open \n",port)
		} else {
			fmt.Printf("Domain : %s, status : %d \n", domain, resp.StatusCode)
			defer resp.Body.Close()
		}
	}
}

func main() {
	check_ports("http://localhost")
}
