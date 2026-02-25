package display

import "fmt"

const Toolname = `
 _   __      _               _
| | / /     | |             | |
| |/ / _   _| |__   ___  ___| |_ _ __ ___  _   _  ___ _ __
|    \| | | | '_ \ / _ \/ __| __| '__/ _ \| | | |/ _ \ '__| v0.3.0
| |\  \ |_| | |_) |  __/\__ \ |_| | | (_) | |_| |  __/ |
\_| \_/\__,_|_.__/ \___||___/\__|_|  \___/ \__, |\___|_|
                                            __/ |
                                           |___/ `

const Author = `			                        By Rolix`

const Split = `
--------------------------------------------------------`

func PrintBanner() {
	fmt.Println("\x1b[1;36m" + Toolname + "\x1b[0m")
	fmt.Println(Author)
	fmt.Println(Split)
}
