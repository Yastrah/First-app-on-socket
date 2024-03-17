import socket
import threading
import logging
import configparser
import sys
from time import sleep

config = configparser.ConfigParser()
config.read('config.ini')


def receive_messages(conn):
    while True:
        try:
            message = conn.recv(int(config['settings']['messageLen'])).decode('utf-8')
            if not message: break
            logging.debug(f"From server received: {message}")

        except socket.error as e:
            logging.error(f"Could not read message: {e}")
            break  

    logging.info("Connection closed by server. Exiting...")


def main():
    file_log = logging.FileHandler(filename="client-logs.log", encoding="utf-8")
    logging.basicConfig(format="%(asctime)s [%(levelname)s]: %(message)s", level="DEBUG",
                        datefmt="%Y-%m-%d %H:%M:%S", handlers=[file_log])

    # подключение к сокету
    conn = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        conn.connect((config['server']['ip'], int(config['server']['port'])))
    except socket.error as e:
        logging.error(f"Could not connect: {e}")
        sys.exit(0)

    logging.info(f"Connected to server <{config['server']['ip']}>")
    
    # запуск потока для чтения данных из сокета
    receive_thread = threading.Thread(target=receive_messages, args=(conn,))
    receive_thread.start()
    
    sleep(int(config['settings']['greetingTime']))

    # отправка сообщения
    info_string = "Страхов Ярослав Константинович. М3О-109Б-23"
    try:
        conn.sendall(info_string.encode('utf-8'))
        logging.debug(f"To server sent: {info_string}")
    except socket.error as e:
        logging.error(f"Connection is inactive: {e}")
        sys.exit(0)


if __name__ == "__main__":
    main()
