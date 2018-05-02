package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"Redis"
)

//退出控制变量
var exit = false

func main() {
	fmt.Println("客户端启动...")
	client,err := Redis.NewClient("127.0.0.1:6379")
	if err != nil {
		log.Panic(err)
	}

	defer client.Close()

	consoleReader := bufio.NewScanner(os.Stdin)
	for !exit && consoleReader.Scan() {
		cmd := consoleReader.Text()
		if strings.ToLower(cmd)=="exit" {
			exit=true
		}else{
			 client.Send(consoleReader.Text())
		}
	}
}


