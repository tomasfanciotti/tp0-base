from socket import socket
import logging

# Config
HEADER_LENGHT = 4
MAX_PACKET_SIZE = 1024

# Encoders and decoders

def encode(data):
    if type(data) == int:
        return str(int).encode()

    if type(data) == str:
        return data.encode()

    if type(data) == list:
        return ("#".join(data) + "#").encode()


def decode(data):
    decoded = data.decode()
    if decoded.isdigit():
        return int(decoded)

    if decoded.find("#") > 0:
        return decoded.split("#")[:-1]

    return decoded


# Application layer data packet, based on TLV format: Type - Lenght - Value

class Packet:

    def __init__(self, opcode: int = None, data_lenght: int = None, data: bytes = None):
        self.opcode = opcode
        self.data_lenght = data_lenght
        self.data = data

    @classmethod
    def new(cls, opcode, data):
        encoded_data = encode(data)
        return Packet(opcode, len(encoded_data), encoded_data)


# Upper Layer

def receive(s: socket):

    # Receive header
    read_bytes = __receive(s, HEADER_LENGHT)
    opcode, data_lenght = int(read_bytes.decode()[0]), int(read_bytes.decode()[1:])

    # Receive data
    to_read = min(data_lenght, MAX_PACKET_SIZE)
    data = __receive(s, to_read)

    while len(data) < data_lenght:

        to_read = min(data_lenght-len(data), MAX_PACKET_SIZE)
        partial_data = __receive(s, to_read)

        if len(partial_data) == 0:
            raise Exception("Che flaco te quedaste corto mepa")

        data += partial_data

    return Packet(opcode, data_lenght, data)


def send(s: socket, packet: Packet):

    # send header
    encoded_header = bytes("{}{:03d}".format(packet.opcode, packet.data_lenght), encoding="utf-8")
    sent_bytes = __send(s, encoded_header)

    # send data
    i, offset = 0, 0
    total_sent = 0

    while total_sent < packet.data_lenght:

        i, offset = i + offset, min(packet.data_lenght - total_sent, MAX_PACKET_SIZE)
        sent_bytes = __send(s, packet.data[i: i+offset])

        if sent_bytes == 0:
            raise Exception("Che flaco te quedaste pagando mepa")

        total_sent += sent_bytes

    return True

# Lower Layer

def __receive(s: socket, total_bytes: int):

    buffer = b''
    actual_read = b''

    while len(buffer) < total_bytes:

        actual_read += s.recv(total_bytes-len(buffer))

        if len(actual_read) == 0:
            # EndOfFile
            break

        if type(actual_read) is int and actual_read < 0:
            # Error
            break

        buffer += actual_read

    logging.debug(f'action: __receive | buffer: {buffer}')
    return buffer


def __send(s: socket, buffer: bytes):
    logging.debug(f'action: __send | buffer: {buffer}')
    sent = 0
    while sent < len(buffer):

        actual_sent = s.send(buffer[sent:])

        if actual_sent == 0:
            # EndOfFile
            break

        if actual_sent < 0:
            # Error
            break

        sent += actual_sent

    return sent
