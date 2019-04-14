package org.apache.plc4x.interop.impl;

import org.apache.plc4x.interop.ConnectionHandle;
import org.apache.plc4x.interop.InteropServer;
import org.apache.plc4x.interop.Request;
import org.apache.plc4x.interop.Response;
import org.apache.thrift.TException;
import org.apache.thrift.protocol.TBinaryProtocol;
import org.apache.thrift.protocol.TProtocol;
import org.apache.thrift.transport.TSocket;
import org.apache.thrift.transport.TTransport;

import java.util.Collections;

public class Client {

    public static void main(String[] args) throws TException {
        try (TTransport transport = new TSocket("localhost", 9090)) {

            transport.open();

            TProtocol protocol = new TBinaryProtocol(transport);

            final InteropServer.Client client = new InteropServer.Client(protocol);

//            final ConnectionHandle connection = client.connect("mock:a");
//            final Response result = client.execute(connection, new Request(Collections.singletonMap("field_1", "DB.field.qry")));
//            System.out.println("Got response: " + result);

            for (int i = 1; i <= 100; i++) {
                final ConnectionHandle connection = client.connect("s7://192.168.167.210/0/1");
                final Response result = client.execute(connection, new Request(Collections.singletonMap("field_1", "%M0:USINT")));
                System.out.println("Got response: " + result);

                client.close(connection);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
