package main

import (
	"fmt"
	"net"
)

func main() {
	listener, _ := net.Listen("tcp", "localhost:8080") // открываем слушающий сокет
	for {
		fmt.Println("Waitinf for new connect")
		conn, err := listener.Accept() // принимаем TCP-соединение от клиента и создаем новый сокет
		if err != nil {
			fmt.Println("Conn is nil!")
			continue
		}
		go handleClient(conn) // обрабатываем запросы клиента в отдельной го-рутине
	}
}

func handleClient(conn net.Conn) {
	fmt.Println("Handle is started")
	defer conn.Close() // закрываем сокет при выходе из функции

	buf := make([]byte, 32)                          // буфер для чтения клиентских данных
	conn.Write([]byte("Hello, what's your name?\n")) // пишем в сокет
	for {
		// fmt.Println("Second For iteration")

		readLen, err := conn.Read(buf) // читаем из сокета, тут сидит ждет подключения

		if err != nil {
			fmt.Println("ERROR", err)
			//  ftm.Println(err.type())
			// fmt.Printf("%T", err)
			break
		}
		
		fmt.Println(string(buf[:readLen]))
		
		answer := make([]byte, 32)
		for i := readLen - 1; i >= 0; i-- {
			answer[readLen - i - 1] = buf[i]
		}

		conn.Write(append([]byte("Echo: "), answer[:readLen]...)) // пишем в сокет
	}
}
