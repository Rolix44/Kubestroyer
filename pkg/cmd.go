package pkg

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/Rolix44/Kubestroyer/utils"
)

func Execute() {

	_, err := os.Stat(utils.Target)
	if err != nil {
		utils.Targets = strings.Split(utils.Target, ",")
	} else {
		file, err := os.Open(utils.Target)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		sc := bufio.NewScanner(file)

		for sc.Scan() {
			utils.Targets = append(utils.Targets, sc.Text())
		}
		if err := sc.Err(); err != nil {
			log.Fatal(err)
		}
	}

	for _, target := range utils.Targets {
		if !utils.AnonRce {
			check_ports(target)
		}

		if utils.AnonRce {
			pods := parse_pod(target)
			anon_rce(pods, target)
		}
	}
}
