from thrift.protocol import TBinaryProtocol
from thrift.transport import TSocket, TTransport
from thrift.transport.TTransport import TTransportException

from org.apache.plc4x.interop.InteropServer import Client, PlcException
from org.apache.plc4x.interop.ttypes import Request

transport = TSocket.TSocket('localhost', 9090)
transport = TTransport.TBufferedTransport(transport)
protocol = TBinaryProtocol.TBinaryProtocol(transport)

client = Client(protocol)

try:
    transport.open()
    connect = None
    try:
        connect = client.connect("mock:a")

        res = client.execute(connect, Request(fields={"field1": "asdf"}))

        print(res)
    finally:
        if connect is not None:
            client.close(connect)
except TTransportException:
    print("Unable to connect to the Interop Server, is the Server really running???")
finally:
    transport.close()
