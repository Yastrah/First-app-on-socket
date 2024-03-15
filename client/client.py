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
            if not message:
                break
        except socket.error as e:
            logging.error(f"Could not recieve message: {e}")
        logging.debug(f"Recieve message: {message}")

def send_messages(conn):
    while True:
        message = input("Write message: ")
        try:
            conn.send(message.encode('utf-8'))
            logging.debug(f"Message sent: {message}")
        except socket.error as e:
            logging.error(f"Connection is inactive: {e}")
            logging.debug("Exiting...")
            sys.exit(0)

def main():
    file_log = logging.FileHandler(filename="test-logs.log", encoding="utf-8")
    # "%(asctime)s [%(levelname)s] %(lineno)d, %(funcName)s: %(message)s",
    logging.basicConfig(format="%(asctime)s [%(levelname)s]: %(message)s", level="DEBUG",
                        datefmt="%Y-%m-%d %H:%M:%S", handlers=[file_log])


    conn = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        conn.connect((config['server']['ip'], int(config['server']['port'])))
    except socket.error as e:
        logging.error(f"Could not connect: {e}")
        logging.debug("Exiting...")
        sys.exit(0)

    logging.debug("Connected!")
    
    receive_thread = threading.Thread(target=receive_messages, args=(conn,))
    receive_thread.start()
    
    sleep(int(config['settings']['greetingTime']))

    info_string = "Страхов Ярослав Константинович. М3О-109Б-23"
    conn.sendall(info_string.encode('utf-8'))
    logging.debug(f"Info sent: {info_string}")
    
    send_messages(conn)
    
if __name__ == "__main__":
    main()
