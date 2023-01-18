package main

import (
	"fmt"

	"github.com/Rolix44/Kubestroyer/pkg"
	"github.com/Rolix44/Kubestroyer/utils"
)

func main() {

	fmt.Println("\x1b[1;36m" + utils.Toolname + "\x1b[0m")
	fmt.Println(utils.Author)
	fmt.Println(utils.Split)

	utils.Config()
	pkg.Execute()

}
