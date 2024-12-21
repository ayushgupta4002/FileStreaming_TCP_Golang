package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type File struct{}

func main() {
	fmt.Println("Hello World")
	go func() {
		time.Sleep(4 * time.Second)
		TcpSend(1000)
	}()
	fs := &File{}
	fs.TcpStart()
}

func (fs *File) TcpStart() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error in listening")
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error in accepting")
		}
		go fs.TcpExecute(conn)
	}
}

func TcpSend(size int) {
	// file := make([]byte, size)
	// _, err := io.ReadFull(rand.Reader, file)
	// if err != nil {
	// 	fmt.Println("Error in reading random")
	// }

	file, err := os.Open("main.txt")
	if err != nil {
		fmt.Println("Error in opening file")
		return
	}
	defer file.Close()

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("Error in dialing")
	}
	defer conn.Close()
	// n, err := conn.Write(file)
	// n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	n, err := io.Copy(conn, file)

	if err != nil {
		fmt.Println("Error in writing")
		fmt.Println(err)
	}
	fmt.Printf("sent %d bytes", n)

	return
}

func (fs *File) TcpExecute(conn net.Conn) {
	defer conn.Close()
	buf := new(bytes.Buffer)
	outputFile, err := os.Create("received_file.txt") // Replace with desired output path
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()
	for {
		// n, err := conn.Read(buf)
		n, err := io.CopyN(outputFile, conn, 100)

		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client")
			} else {
				fmt.Println("Error in reading:", err)
			}
			break // Exit the loop on error
		}
		fmt.Println(buf.Bytes())
		fmt.Printf("receiver %d bytes", n)
	}
	//output shall be 1000
	fmt.Printf("Final buffer length: %d bytes\n", buf.Len())

}
