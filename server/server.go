package main

import (
	"log"
	"net"
	"os"
	"time"

	"gopkg.in/ini.v1"
)

type Config struct {
	host       string
	port       string
	messageLen int
	answerTime int
	deadline   int
}

var config Config

func main() {
	config = loadConfig()

	file, err := os.OpenFile("test-logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("Failed to open debug log file:", err)
	}

	log.SetOutput(file)
	log.Println("Starting listening...")

	listener, _ := net.Listen("tcp", config.host+":"+config.port) // открываем слушающий сокет
	for {
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

	// Установка тайм-аута для подключения
	deadline := time.Now().Add(time.Duration(config.deadline) * time.Second)
	err := conn.SetReadDeadline(deadline)
	if err != nil {
		log.Printf("Error setting read deadline: %v", err)
		return
	}

	buf := make([]byte, config.messageLen) // буфер для чтения клиентских данных. Может принять "messageLen" символа за раз

	for {
		readLen, err := conn.Read(buf) // читаем из сокета, тут сидит ждет новые данные

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() { // Соединение закрыто по тайм-ауту
				log.Printf("Connection %s is closed due to timeout!\n", conn.RemoteAddr().String())
			} else { // Другая ошибка
				log.Printf("Connection %s is closed! Error: %v\n", conn.RemoteAddr().String(), err)
			}
			break
		}

		log.Printf("From %s received: %s", conn.RemoteAddr().String(), string(buf[:readLen]))

		reversedMess := make([]byte, config.messageLen)
		for i := readLen - 1; i >= 0; i-- {
			reversedMess[readLen-i-1] = buf[i]
		}
		answer := string(reversedMess[:readLen]) + ". Сервер разработан Страховым Я.K. M3O-109Б-23"

		time.Sleep(time.Duration(config.answerTime) * time.Second) // Симуляция работы сервера
		conn.Write([]byte(answer))

		log.Printf("To %s sent: %s", conn.RemoteAddr().String(), answer)
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
	config.answerTime, _ = section.Key("answerTime").Int()
	config.deadline, _ = section.Key("deadline").Int()

	return config
}
