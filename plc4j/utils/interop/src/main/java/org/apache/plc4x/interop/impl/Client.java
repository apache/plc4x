/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

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
