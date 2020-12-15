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

import subprocess
import time
import warnings
from contextlib import contextmanager

from generated.org.apache.plc4x.interop.InteropServer import Client, PlcException
from thrift.protocol import TBinaryProtocol
from thrift.transport import TSocket, TTransport
from thrift.transport.TTransport import TTransportException

from org.apache.plc4x.PlcConnection import PlcConnection


class PlcDriverManager:
    """
    Constructor, initialize the server
    """
    def __init__(self, embedded_server=True):
        self.embedded_server = embedded_server
        self.interop_proc = None
        # Start the Server in the background
        if embedded_server:
            self.start_server()

        transport = TSocket.TSocket('localhost', 9090)
        self.transport = TTransport.TBufferedTransport(transport)
        self.protocol = TBinaryProtocol.TBinaryProtocol(self.transport)

    def __enter__(self):
        self.open()
        return self

    def __exit__(self, *args):
        self.close()

    def start_server(self):
        self.interop_proc = subprocess.Popen(
            ["java", "-Dlog4j.configurationFile=../lib/log4j2.xml", "-jar", "../lib/interop-server.jar"]
        )
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

    def _get_client(self):
        return Client(self.protocol)

    def get_connection(self, url):
        return PlcConnection(self._get_client(), url)

    @contextmanager
    def connection(self, url):
        """
        Context manager to handle connection.
        """
        conn = None
        try:
            conn = self.get_connection(url)
            yield conn
        except PlcException as e:
            raise Exception(str(e.url))
        finally:
            if conn is not None:
                conn.close()

    def open(self):
        try:
            self.transport.open()
        except TTransportException:
            self.close()  # Handle failure on enter
            raise PlcException("Unable to connect to the Interop Server, is the Server really running?")

    def close(self):
        print("Closing the Interop Server")
        try:
            self.transport.close()
        finally:
            if self.embedded_server:
                self.interop_proc.terminate()
