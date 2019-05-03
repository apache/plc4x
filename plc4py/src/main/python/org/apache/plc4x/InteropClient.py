#  Licensed to the Apache Software Foundation (ASF) under one
#  or more contributor license agreements.  See the NOTICE file
#  distributed with this work for additional information
#  regarding copyright ownership.  The ASF licenses this file
#  to you under the Apache License, Version 2.0 (the
#  "License"); you may not use this file except in compliance
#  with the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing,
#  software distributed under the License is distributed on an
#  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
#  KIND, either express or implied.  See the License for the
#  specific language governing permissions and limitations
#  under the License.

from generated.org.apache.plc4x.interop.InteropServer import Client
from generated.org.apache.plc4x.interop.ttypes import Request
from thrift.protocol import TBinaryProtocol
from thrift.transport import TSocket, TTransport
from thrift.transport.TTransport import TTransportException

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
