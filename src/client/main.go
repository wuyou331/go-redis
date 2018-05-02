package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

//退出控制变量
var exit = false

func main() {
	fmt.Println("客户端启动...")

	client, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Panic(err)
	}

	defer client.Close()
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	reader := bufio.NewReader(client)
	go responseHandle(reader, &waitGroup)

	writer := bufio.NewWriter(client)
	go requestHandle(writer, &waitGroup)

	waitGroup.Wait()
}

func requestHandle(writer *bufio.Writer, wait *sync.WaitGroup) {

	consoleReader := bufio.NewScanner(os.Stdin)
	for !exit && consoleReader.Scan() {
		line := consoleReader.Text()
		var inputs []string
		for _, input := range strings.Split(line, " ") {
			if input != "" {
				inputs = append(inputs, input)
			}
		}
		cmd := packRequest(inputs)
		writer.WriteString(cmd)
		writer.Flush()
	}
	wait.Done()
}

//打包请求命令
func packRequest(inputs []string) string {
	var cmd = fmt.Sprintf("*%d\r\n", len(inputs))
	for _, arg := range inputs {
		cmd = cmd + fmt.Sprintf("$%d\r\n", len(arg))
		cmd = cmd + fmt.Sprintf("%s\r\n", arg)
	}
	return cmd
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
			msg := strings.Trim(string(line[1:]), "\r\n")
			fmt.Println(msg)
			break
		case '-':
			msg := strings.Trim(string(line[1:]), "\r\n")
			log.Print(msg)
			break
		case ':':
			msg := strings.Trim(string(line[1:]), "\r\n")
			fmt.Println(msg)
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
				reader.Discard(msgLength + 2)
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
