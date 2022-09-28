import socket
from dataclasses import dataclass
import time


@dataclass
class Server:
    host: str
    port: int

    def __post_init__(self):
        self._sock = socket.socket(socket.AF_INET,
                                   socket.SOCK_STREAM)
        self._sock.setsockopt(socket.SOL_SOCKET,
                              socket.SO_REUSEADDR, 1)

    def __enter__(self):
        self._sock.bind((self.host, self.port))
        return self

    def __exit__(self, exception_type, exception_value, traceback):
        self._sock.close()

    def listen_for_traffic(self):
        self._sock.listen(5)
        connection, address = self._sock.accept()
        with connection:
            while True:
                message = connection.recv(1024)
                count = connection.send(message)
        self._sock.close()
