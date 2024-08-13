package main

import (
	"flag"
	"fmt"
	"github.com/jacksonopp/monkey/repl"
	"os"
	user2 "os/user"
)

func main() {
	user, err := user2.Current()
	if err != nil {
		panic(err)
	}

	evalFlag := flag.Bool("eval", true, "start in eval mode")
	astFlag := flag.Bool("ast", false, "start in ast mode")
	tokenFlag := flag.Bool("token", false, "start in token mode")

	flag.Parse()

	flags := [3]*bool{evalFlag, astFlag, tokenFlag}

	fmt.Println("astFlag", *astFlag, "tokenFlag", *tokenFlag)

	fmt.Printf("Hello %s! This is the Monkey programming language\n", user.Username)
	if *astFlag {
		fmt.Printf("Starting in ast mode\n")
	} else if *tokenFlag {
		fmt.Printf("Starting in token mode\n")
	} else {
		fmt.Printf("Starting in eval mode\n")
	}
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout, flags)
}
