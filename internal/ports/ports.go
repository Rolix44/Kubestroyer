package ports

import "maps"

var knownPorts = map[int]string{
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

// KnownPorts returns a clone of the base Kubernetes port map.
func KnownPorts() map[int]string {
	return maps.Clone(knownPorts)
}

// WithNodePorts returns the base port map extended with NodePort range 30000–32767.
func WithNodePorts() map[int]string {
	m := KnownPorts()
	for p := 30000; p <= 32767; p++ {
		m[p] = "Possible NodePort service"
	}
	return m
}
