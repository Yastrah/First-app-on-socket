package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const messageLen int = 128

func main() {
	file, err := os.OpenFile("test-logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("Failed to open debug log file:", err)
	}

	debugLog := log.New(file, "[DEBUG] ", log.Ldate|log.Ltime)
	errorLog := log.New(file, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

	// log.SetOutput(file)
	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	debugLog.Println("some log info")
	errorLog.Println("some log info")

	listener, _ := net.Listen("tcp", "192.168.1.106:8080") // открываем слушающий сокет
	for {
		// fmt.Println("Waiting for new connect...")
		conn, err := listener.Accept() // принимаем TCP-соединение от клиента и создаем новый сокет
		if err != nil {
			log.Printf("Conn is nil!")
			continue
		}

		// fmt.Println("New connect accepted")
		log.Printf("Received connection from %s", conn.RemoteAddr().String())

		go handleClient(conn) // обрабатываем запросы клиента в отдельной го-рутине
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close() // закрываем сокет при выходе из функции

	buf := make([]byte, messageLen)                      // буфер для чтения клиентских данных. Может принять 32 символа за раз
	conn.Write([]byte("Hello, what's your name?\n")) // пишем в сокет
	for {
		readLen, err := conn.Read(buf) // читаем из сокета, тут сидит ждет новые данные

		if err != nil {
			log.Printf("Connection %s is closed!", conn.RemoteAddr().String())
			break
		}

		fmt.Println(string(buf[:readLen]))

		answer := make([]byte, messageLen)
		for i := readLen - 1; i >= 0; i-- {
			answer[readLen-i-1] = buf[i]
		}
		str_answer := "Answer: <" + string(answer[:readLen]) + "> - Server developed by Yastrah"

		// conn.Write(append([]byte("Echo: "), answer[:readLen]...)) // пишем в сокет
		conn.Write([]byte(str_answer))
	}
}
