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
package org.apache.plc4x.java.utils.pcapsockets.netty;

import java.io.File;
import java.net.InetAddress;
import java.net.SocketAddress;

public class PcapSocketAddress extends SocketAddress {

    private static final long serialVersionUID = 1L;

    public static final int ALL_PORTS = -1;
    public static final int ALL_PROTOCOLS = -1;

    private final File pcapFile;
    private final InetAddress address;
    private final int port;
    private final int protocolId;

    public PcapSocketAddress(File pcapFile, InetAddress address, int port, int protocolId) {
        this.pcapFile = pcapFile;
        this.address = address;
        this.port = port;
        this.protocolId = protocolId;
    }

    public File getPcapFile() {
        return pcapFile;
    }

    public InetAddress getAddress() {
        return address;
    }

    public int getPort() {
        return port;
    }

    public int getProtocolId() {
        return protocolId;
    }

}
