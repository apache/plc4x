import time
import unittest
from subprocess import Popen
from unittest import TestCase

from org.apache.plc4x.PlcDriverManager import PlcDriverManager
from org.apache.plc4x.interop.InteropServer import Client, Request, PlcException

if __name__ == '__main__':
    unittest.main()


class TestPlcDriverManager(TestCase):

    def test_callLib(self):
        try:
            proc = Popen(["java", "-jar", "../lib/apache-plc4x-incubating-0.4.0-SNAPSHOT-jar-with-dependencies.jar"])
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

    def test_withRealPLC_forDebug(self):
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

    def test_withRealPLC(self):
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