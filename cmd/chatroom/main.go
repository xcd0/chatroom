package main

import (
	"flag"
	"log"
	"os"

	chat "github.com/xcd0/chatroom"
)

func main() {
	err := chat.Run(os.Args[1:], os.Stdout, os.Stderr)
	if err != nil && err != flag.ErrHelp {
		log.Println(err)
		exitCode := 1
		if ecoder, ok := err.(interface{ ExitCode() int }); ok {
			exitCode = ecoder.ExitCode() // errの中から終了コードが取り出せる場合はその終了コードを返す。
		}
		os.Exit(exitCode)
	}
}
