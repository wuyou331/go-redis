package Redis

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type client struct {
	conn   net.Conn
	writer *bufio.Writer
	reader *bufio.Reader
}
type Client interface {
	Close()
	Send(cmd string) ReplyMessage
}

//New 根据Ip地址和端口创建一个RedisClient
func NewClient(ipAndPort string) (Client, error) {
	client := client{}
	conn, err := net.Dial("tcp", ipAndPort)

	client.conn = conn
	client.writer = bufio.NewWriter(conn)
	client.reader = bufio.NewReader(conn)

	return &client, err
}

func (client *client) Send(cmd string) ReplyMessage {
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

	return client.waitReplyMessage()
}

func (client *client) Close() {
	client.conn.Close()
}

func (client *client) waitReplyMessage() ReplyMessage {
	line, err := client.reader.ReadBytes('\n')
	if err != nil {
		log.Panic(err)
	}

	cmdType := line[0]
	reply := replyMessage{}
	switch cmdType {
	case '+':
		reply.replyType = SingleLine
		msg := strings.Trim(string(line[1:]), "\r\n")
		reply.message = msg
	case ':':
		reply.replyType = Integer
		msg := strings.Trim(string(line[1:]), "\r\n")
		reply.message = msg
	case '-':
		reply.replyType = Error
		msg := strings.Trim(string(line[1:]), "\r\n")
		reply.message = msg
	case '$':
		reply.replyType = Bulk
		msgLength, _ := strconv.Atoi(strings.Trim(string(line[1:]), "\r\n"))
		block, err := client.reader.Peek(msgLength)
		if err != nil {
			log.Panic(err)
		}
		client.reader.Discard(msgLength + len([]byte("\r\n")))
		msg := string(block)
		reply.message = msg
	default:
		reply.replyType = Unknown
	}
	return &reply
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
