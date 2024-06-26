package utils

func init() {
	addKnownServicePort()
}

var KnownPorts = map[int]string{
	443:   "Kubernetes API port",
	2379:  "Etcd",
	6666:  "Etcd",
	4194:  "cAdvisor for containers metrics",
	6443:  "Kubernetes API port",
	8443:  "Minikube API port",
	8080:  "Insecure API port",
	10250: "Kubelet API anonymous port",
	10255: "Kubelet API read only",
	10256: "Kube proxy health check server",
	9099:  "Calico health check server",
	6782:  "Weave metrics/endpoints",
	6783:  "Weave metrics/endpoints",
	6784:  "Weave metrics/endpoints",
	44134: "Tiller service listening",
}

func addKnownServicePort() {
	for port := 30000; port <= 32767; port++ {
		KnownPorts[port] = "Possible NodePort service"
	}
}
