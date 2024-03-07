package main

import (
	"log"
	"net"
	"os"
	"gopkg.in/ini.v1"
)

type Config struct {
	host       string
	port       string
	messageLen int
}

var config Config

func main() {
	config = loadConfig()

	file, err := os.OpenFile("test-logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("Failed to open debug log file:", err)
	}

	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println("Starting listening...")

	listener, _ := net.Listen("tcp", config.host+":"+config.port) // открываем слушающий сокет
	for {
		// fmt.Println("Waiting for new connect...")
		conn, err := listener.Accept() // принимаем TCP-соединение от клиента и создаем новый сокет
		if err != nil {
			log.Printf("Conn is nil!")
			continue
		}

		log.Printf("Received connection from %s", conn.RemoteAddr().String())

		go handleClient(conn) // обрабатываем запросы клиента в отдельной го-рутине
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close() // закрываем сокет при выходе из функции

	buf := make([]byte, config.messageLen)           // буфер для чтения клиентских данных. Может принять 32 символа за раз
	conn.Write([]byte("Hello, what's your name?\n")) // пишем в сокет
	for {
		readLen, err := conn.Read(buf) // читаем из сокета, тут сидит ждет новые данные

		if err != nil {
			log.Printf("Connection %s is closed!", conn.RemoteAddr().String())
			break
		}

		log.Printf("From %s received: %s", conn.RemoteAddr().String(), string(buf[:readLen]))

		answer := make([]byte, config.messageLen)
		for i := readLen - 1; i >= 0; i-- {
			answer[readLen-i-1] = buf[i]
		}
		str_answer := "Answer: <" + string(answer[:readLen]) + "> - Server developed by Yastrah"

		// conn.Write(append([]byte("Echo: "), answer[:readLen]...)) // пишем в сокет
		conn.Write([]byte(str_answer))

		log.Printf("To %s sent: %s", conn.RemoteAddr().String(), string(str_answer))
	}
}

func loadConfig() Config {
	inidata, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal("Fail to read file: ", err)
	}
	section := inidata.Section("settings")

	var config Config
	config.host = section.Key("host").String()
	config.port = section.Key("port").String()
	config.messageLen, _ = section.Key("messageLen").Int()

	return config
}
