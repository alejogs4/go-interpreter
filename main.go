package main

import (
	"fmt"
	"os"
	osUser "os/user"

	"go-interpreter.com/m/repl"
)

func main() {
	user, err := osUser.Current()
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Hello %s feel free to type commands", user.Username))
	repl.Start(os.Stdin, os.Stdout)
}
