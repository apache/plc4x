/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.sandbox.java.dynamic.io;

import org.apache.plc4x.sandbox.java.dynamic.exceptions.DynamicDriverException;

import java.net.DatagramSocket;
import java.net.SocketException;

public class UdpProtocolIO implements ProtocolIO {

    private DatagramSocket socket;

    public UdpProtocolIO(String host, int port) throws DynamicDriverException {
        try {
            socket = new DatagramSocket(port);
            //socket.receive();
        } catch (SocketException e) {
            throw new DynamicDriverException("Error creating UDP Socket", e);
        }
    }

    @Override
    public void send(byte[] data) {

    }

    @Override
    public byte[] receive() {
        return new byte[0];
    }

}
