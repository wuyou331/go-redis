package Redis

import (
	"net"
	"fmt"
	"strings"
	"bufio"
	"log"
	"strconv"
)

type client struct {
	conn net.Conn
	writer *bufio.Writer
	reader *bufio.Reader
}
type Client interface{
	Close()
	Send(cmd string)
}

//New 根据Ip地址和端口创建一个RedisClient
func NewClient(ipAndPort string) (Client,error){
	client := client{}
	conn, err := net.Dial("tcp", ipAndPort)

	client.conn=conn
	client.writer=bufio.NewWriter(conn)
	client.reader=bufio.NewReader(conn)

	return  &client,err
}

func (client *client) Send(cmd string) {
	var cmdArr []string
	//移除空白的项目
	for _, input := range strings.Split(cmd, " ") {
		if input != "" {
			cmdArr = append(cmdArr, input)
		}
	}
	cmd = packRequest(cmdArr)
	client.writer.WriteString(cmd)
	client.writer.Flush()

	line, err := client.reader.ReadBytes('\n')
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
			block, err := client.reader.Peek(msgLength)
			if err != nil {
				log.Panic(err)
			}
			client.reader.Discard(msgLength +len([]byte("\r\n")))
			msg := string(block)
			fmt.Println(msg)
		}
		break
	default:
		break
	}
}

func (client *client) Close() {
	client.conn.Close()
}

//packRequest 方法用于将用户输入的命令打包成redis的请求格式
func packRequest(cmds []string) string {
	var cmd = fmt.Sprintf("*%d\r\n", len(cmds))
	for _, arg := range cmds {
		cmd = cmd + fmt.Sprintf("$%d\r\n", len(arg))
		cmd = cmd + fmt.Sprintf("%s\r\n", arg)
	}
	return cmd
}
