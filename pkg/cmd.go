package pkg

import "github.com/Rolix44/Kubestroyer/utils"

func Execute() {

	if !utils.AnonRce {
		check_ports(utils.Target)
	}

	if utils.AnonRce {
		pods := parse_pod(utils.Target)
		anon_rce(pods)
	}

}
