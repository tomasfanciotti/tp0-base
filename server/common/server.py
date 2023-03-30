import socket
import logging
import signal
from .messaging_protocol import Packet, decode, receive, send
from .national_lottery import *

# OpCodes
OP_CODE_ZERO = 0
OP_CODE_REGISTER = 1
OP_CODE_REGISTER_ACK = 2
OP_CODE_REGISTER_BATCH = 3
OP_CODE_ERROR = 4


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.listening = True
        signal.signal(signal.SIGTERM, self.__stop_listening)

    def __stop_listening(self, *args):
        self.listening = False
        self._server_socket.close()

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while self.listening:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
            except OSError as e:
                logging.error(f"action: accept_connections | result: fail | error: {e}")

        logging.info(f"action: run | result: succes | msg: server shutting down ")

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            while True:
                packet = receive(client_sock)
                addr = client_sock.getpeername()
                logging.debug(
                    f'action: receive_message | result: success | ip: {addr[0]} | msg: {packet.data.decode()}')

                if packet.opcode == OP_CODE_REGISTER:

                    argv = decode(packet.data)
                    dni, number = register_bet(argv)
                    if dni and number:
                        logging.info(f'action: apuesta_almacenada | result: success | dni: ${dni} | numero: ${number}')
                        response = Packet.new(OP_CODE_REGISTER_ACK, "")
                        send(client_sock, response)

                    else:
                        logging.error(f'action: apuesta_almacenada | result: fail | ip: ${addr[0]} | args: ${argv}')
                        response = Packet.new(OP_CODE_ERROR, "Error al procesar la apuesta enviada")
                        send(client_sock, response)

                elif packet.opcode == OP_CODE_REGISTER_BATCH:

                    logging.debug(
                        f'action: procesar_batch | result: in_progress | ip: ${addr[0]} | arg len: ${len(packet.data)}')

                    argv = decode(packet.data)
                    result = register_batch(argv)
                    if result:
                        response = Packet.new(OP_CODE_REGISTER_ACK, "")
                        send(client_sock, response)

                        logging.info(f'action: procesar_batch | result: success | ip: {addr[0]} | stored: {result}')
                    else:
                        response = Packet.new(OP_CODE_ERROR, "Error al procesar el batch enviado")
                        send(client_sock, response)

                        logging.error(f'action: procesar_batch | result: fail | ip: {addr[0]} | stored: {0}')

                elif packet.opcode == OP_CODE_ZERO:
                    logging.info(f'action: disconnected | result: success | ip: {addr[0]} ')
                    break

        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
