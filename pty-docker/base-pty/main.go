package main

import (
	"os/exec"
	"time"

	"github.com/creack/pty"
)

func main() {
	c := exec.Command("python")
	f, _ := pty.Start(c)
	defer f.Close()
	time.Sleep(time.Second)

	f.Write([]byte("print('hello')\n"))
	f.Write([]byte{4}) // EOT
	// go func() {
	// 	f.Write([]byte("ls\n"))
	// 	f.Write([]byte{4}) // EOT
	// }()

}
