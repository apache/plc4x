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

import time
import unittest
from subprocess import Popen
from unittest import TestCase

from generated.org.apache.plc4x.interop.InteropServer import Request, PlcException

from org.apache.plc4x.PlcDriverManager import PlcDriverManager

if __name__ == '__main__':
    unittest.main()


class TestPlcDriverManager(TestCase):

    def test_callLib(self):
        try:
            proc = Popen(["java", "-jar", "../lib/interop-server.jar"])
            print("Started server under pid " + str(proc.pid))
            time.sleep(5)
            poll = proc.poll()
            if poll is None:
                print("Still running...")
            else:
                print("Unable to start the Interop Server...")
        except:
            print("Unable to start proces")
        finally:
            proc.terminate()
            time.sleep(1)

    def test_startAndStopServer(self):
        try:
            manager = PlcDriverManager()

            ## Do some magic here
            client = manager._get_client()
            connection = client.connect("mock:a")

            result = client.execute(connection, Request(fields={"field1": "asdf"}))

            print(result)

            client.close(connection)
        finally:
            manager.close()
            time.sleep(1)

    def test_withPlcConnection(self):
        try:
            manager = PlcDriverManager()

            try:
                connection = manager.get_connection("mock:a")
                result = connection.execute(Request(fields={"field1": "asdf"}))
                print(result)
            finally:
                connection.close()
        finally:
            manager.close()
            time.sleep(1)

    def test_withCompleteAPI(self):
        try:
            manager = PlcDriverManager()

            connection = None
            try:
                connection = manager.get_connection("mock:a")

                result = connection.execute(Request(fields={"field1": "asdf"}))
                print("Response Code is " + str(result.get_field("field1").get_response_code()))
                # We now that we want to get an int...
                print("Response Value is " + str(result.get_field("field1").get_int_value()))

                ## Assert the value
                self.assertEqual(result.get_field("field1").get_int_value(), 100)

            except PlcException as e:
                raise Exception(str(e.url))
            finally:
                if connection is not None:
                    connection.close()
        finally:
            manager.close()
            time.sleep(1)

    # This i a manual test, which needs a running server to work
    def mt_withRealPLC_forDebug(self):
        try:
            manager = PlcDriverManager(embedded_server=False)

            connection = None
            try:
                connection = manager.get_connection("s7://192.168.167.210/0/1")
                result = connection.execute(Request(fields={"field1": "%M0:USINT"}))
                print(result)
            except PlcException as e:
                raise Exception(str(e.url))
            finally:
                if connection is not None:
                    connection.close()
        finally:
            manager.close()

    # This is an integration test, which needs a S7 PLC
    def it_withRealPLC(self):
        try:
            manager = PlcDriverManager()

            connection = None
            try:
                connection = manager.get_connection("s7://192.168.167.210/0/1")
                for _ in range(100):
                    result = connection.execute(Request(fields={"field1": "%M0:USINT"}))
                    print("Response Code is " + str(result.get_field("field1").get_response_code()))
                    # We now that we want to get an int...
                    print("Response Value is " + str(result.get_field("field1").get_int_value()))

            except PlcException as e:
                raise Exception(str(e.url))
            finally:
                if connection is not None:
                    connection.close()
        finally:
            manager.close()