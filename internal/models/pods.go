package models

// Container holds a pod container's name and image.
type Container struct {
	Name      string   `json:"name"`
	Image     string   `json:"image"`
	Resources struct{} `json:"resources"`
}

// PodSpec holds the list of containers in a pod.
type PodSpec struct {
	Containers []Container `json:"containers"`
}

// PodMetadata holds identifying metadata for a pod.
type PodMetadata struct {
	Name              string      `json:"name"`
	Namespace         string      `json:"namespace"`
	UID               string      `json:"uid"`
	CreationTimestamp interface{} `json:"creationTimestamp"`
}

// PodItem represents a single running pod.
type PodItem struct {
	Metadata PodMetadata `json:"metadata"`
	Spec     PodSpec     `json:"spec"`
	Status   struct{}    `json:"status"`
}

// RunningPods is the response body from /runningpods/.
type RunningPods struct {
	Kind       string    `json:"kind"`
	APIVersion string    `json:"apiVersion"`
	Metadata   struct{}  `json:"metadata"`
	Items      []PodItem `json:"items"`
}
