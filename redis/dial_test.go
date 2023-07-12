package redis

import (
	"bufio"
	"fmt"
	"net"
	"testing"
)

func Test_Redis_Dial(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	command1 := "SET key1 value1\r\n"
	command2 := "GET key1\r\n"

	_, err = writer.WriteString(command1)
	if err != nil {
		panic(err)
		return
	}
	err = writer.Flush()
	if err != nil {
		panic(err)
		return
	}

	response1, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
		return
	}

	_, err = writer.WriteString(command2)
	if err != nil {
		panic(err)
		return
	}
	err = writer.Flush()
	if err != nil {
		panic(err)
		return
	}

	response2, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
		return
	}

	fmt.Println("response1:", response1)
	fmt.Println("response2:", response2)
}
