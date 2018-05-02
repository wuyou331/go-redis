package Redis

import (
	"net"
	"fmt"
	"strings"
	"bufio"
)

type redisClient struct {
	conn net.Conn
	writer *bufio.Writer
	reader *bufio.Reader
}
type RedisClient interface{
	Close()
	Send(cmd string)
}

//New 根据Ip地址和端口创建一个RedisClient
func New(ipAndPort string) (RedisClient,error){
	client := redisClient{}
	conn, err := net.Dial("tcp", ipAndPort)

	client.conn=conn
	client.writer=bufio.NewWriter(conn)
	client.reader=bufio.NewReader(conn)

	return  &client,err
}

func (client *redisClient) Send(cmd string) {
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


}

func (client *redisClient) Close() {
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
