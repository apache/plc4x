import subprocess
import time
import warnings

from thrift.protocol import TBinaryProtocol
from thrift.transport import TSocket, TTransport
from thrift.transport.TTransport import TTransportException

from org.apache.plc4x.PlcConnection import PlcConnection
from org.apache.plc4x.interop.InteropServer import Client, PlcException


class PlcDriverManager:

    """
    constructor, initialize the server
    """
    def __init__(self, embedded_server=True):
        self.embedded_server = embedded_server
        # Start the Server in the background
        if embedded_server:
            self.interop_proc = subprocess.Popen(
                ["java", "-Dlog4j.configurationFile=../lib/log4j2.xml",
                 "-jar", "../lib/interop-server.jar"])
            try:
                print("Started server under pid " + str(self.interop_proc.pid))
            except:
                print("Encountered an Exception while starting Interop Server")
                raise PlcException("Unable to start the Interop Server!")

            time.sleep(2)
            poll = self.interop_proc.poll()
            if poll is None:
                print("Sucesfully started the Interop Server...")
            else:
                warnings.warn("Interop Server died after starting up...")
                raise PlcException(
                    "Unable to start the Interop Server. Is another Server still running under the same port?")

        self.transport = TSocket.TSocket('localhost', 9090)
        self.transport = TTransport.TBufferedTransport(self.transport)

        self.protocol = TBinaryProtocol.TBinaryProtocol(self.transport)

        try:
            self.transport.open()
        except TTransportException:
            raise PlcException("Unable to connect to the Interop Server, is the Server really running???")

    def _get_client(self):
        return Client(self.protocol)

    def get_connection(self, url):
        return PlcConnection(self._get_client(), url)

    def close(self):
        print("Closing the Interop Server")
        try:
            self.transport.close()
        finally:
            if self.embedded_server:
                self.interop_proc.terminate()
