package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

        func MessageInfom(message string) {
	  fmt.Println("Inform: " + message)
        }

func main() {
	ServerAddr, err := net.ResolveUDPAddr("udp4", "172.17.0.4:10001")
	CheckError(err)

        MessageInfom( "After ResolveUDPAddr")

	//myAddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:10002")
	//CheckError(err)

	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	CheckError(err)

       MessageInfom( "After DialUPD ")

	defer Conn.Close()
	i := 0
	for {
		msg := "Pump" + strconv.Itoa(i)
		i++
		buf := []byte(msg)
		_, err := Conn.Write(buf)
		if err != nil {
			fmt.Println(msg, err)
		}
		time.Sleep(time.Second * 1)
	}
}
