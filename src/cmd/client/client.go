package main

import (
	"Redis"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

//退出控制变量
var exit = false

func main() {
	fmt.Println("客户端启动...")
	client, err := Redis.NewClient("127.0.0.1:6379")
	if err != nil {
		log.Panic(err)
	}

	defer client.Close()

	consoleReader := bufio.NewScanner(os.Stdin)
	for !exit && consoleReader.Scan() {
		cmd := consoleReader.Text()
		if strings.ToLower(cmd) == "exit" {
			exit = true
		} else {
			reply := client.Send(consoleReader.Text())
			printReply(reply)
		}
	}
}

func printReply(reply Redis.ReplyMessage) {

	switch reply.GetType() {
	case Redis.SingleLine:
		fallthrough
	case Redis.Integer:
		fallthrough
	case Redis.Bulk:
		fmt.Println(reply.GetMessage())
	case Redis.Error:
		log.Println(reply.GetMessage())
	case Redis.Unknown:
		log.Println("Unknown")
	}
}
