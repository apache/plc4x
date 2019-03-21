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

import java.io.IOException;
import java.net.Socket;

public class TcpProtocolIO implements ProtocolIO {

    private Socket socket;

    public TcpProtocolIO(String host, int port) throws DynamicDriverException {
        try {
            this.socket = new Socket(host, port);
        } catch (IOException e) {
            throw new DynamicDriverException("Error creating TCP Socket", e);
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
