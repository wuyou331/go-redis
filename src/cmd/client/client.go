package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"Redis/Client"
)

//退出控制变量
var exit = false

func main() {
	fmt.Println("客户端启动...")
	client,err := Client.New("127.0.0.1:6379")
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
			go client.Send(consoleReader.Text())
		}
	}
}



func responseHandle(reader *bufio.Reader, wait *sync.WaitGroup) {
	for !exit {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Panic(err)
		}
		cmdType := line[0]
		switch cmdType {
		case '+':
		case ':':
			msg := strings.Trim(string(line[1:]), "\r\n")
			fmt.Println(msg)
			break
		case '-':
			msg := strings.Trim(string(line[1:]), "\r\n")
			log.Print(msg)
			break
		case '$':
			msgLength, _ := strconv.Atoi(strings.Trim(string(line[1:]), "\r\n"))
			if msgLength == -1 {
				fmt.Println("nil")
			} else {
				block, err := reader.Peek(msgLength)
				if err != nil {
					log.Panic(err)
				}
				reader.Discard(msgLength +len([]byte("\r\n")))
				msg := string(block)
				fmt.Println(msg)
			}
			break
		default:
			break
		}
	}
	wait.Done()
}
