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
package org.apache.plc4x.java.s7.connection;

import org.apache.plc4x.java.connection.PlcConnection;
import org.apache.plc4x.java.exceptions.PlcException;
import org.apache.plc4x.java.exceptions.PlcIoException;
import org.apache.plc4x.java.exceptions.PlcProtocolException;
import org.apache.plc4x.java.s7.messages.ConnectionRequest;
import org.apache.plc4x.java.s7.messages.ConnectionResponse;
import org.apache.plc4x.java.s7.messages.Message;
import org.apache.plc4x.java.s7.messages.NegotiatePduRequest;
import org.apache.plc4x.java.s7.messages.s7.S7Parser;
import org.apache.plc4x.java.s7.messages.s7.messages.S7ResponseMessage;
import org.apache.plc4x.java.s7.messages.s7.params.SetupCommunicationParameter;
import org.apache.plc4x.java.s7.utils.PcapngUtils;

import java.io.*;
import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.UnknownHostException;
import java.nio.ByteBuffer;
import java.nio.channels.ClosedChannelException;
import java.nio.channels.SocketChannel;

public class S7PlcConnection implements PlcConnection {

    private static final int ISO_ON_TCP_PORT = 102;

    private final String hostName;
    private final int rack;
    private final int slot;

    private int pduSize;

    private SocketChannel channel;

    public S7PlcConnection(String hostName, int rack, int slot) {
        this.hostName = hostName;
        this.rack = rack;
        this.slot = slot;
        this.pduSize = 1024;
    }

    public String getHostName() {
        return hostName;
    }

    public int getRack() {
        return rack;
    }

    public int getSlot() {
        return slot;
    }

    public int getPduSize() {
        return pduSize;
    }

    @Override
    public void connect() throws PlcException {
        try {
            // Generate some connection id.
            short connectionId = (short) (Math.random() % Short.MAX_VALUE);

            InetAddress serverIPAddress = InetAddress.getByName(hostName);
            InetSocketAddress serverAddress = new InetSocketAddress(serverIPAddress, ISO_ON_TCP_PORT);
            channel = SocketChannel.open();
            channel.connect(serverAddress);

            // Create a connection request and send it to the PLC.
            System.out.println("Connecting to S7 PLC on " + hostName + " ...");
            ConnectionRequest connectionRequest = new ConnectionRequest(connectionId, (byte) rack, (byte) slot);
            ByteArrayOutputStream baos = new ByteArrayOutputStream(pduSize);
            DataOutputStream dos = new DataOutputStream(baos);
            connectionRequest.serialize(dos);
            channel.write(ByteBuffer.wrap(baos.toByteArray()));

            S7Parser parser = new S7Parser();
            // Read the response.
            ByteBuffer ackBuffer = ByteBuffer.allocate(pduSize);
            channel.read(ackBuffer);
            DataInputStream dis = new DataInputStream(new ByteArrayInputStream(ackBuffer.array()));
            Message message = parser.parse(dis);
            if(!(message instanceof ConnectionResponse)) {
                throw new PlcProtocolException("Expecting Connection Response");
            }
            System.out.println("Connected");

            // Negotiate the PDU Length
            System.out.println("Negotiating PDU (Suggested PDU of 960) ...");
            baos.reset();
            NegotiatePduRequest negotiatePduRequest = new NegotiatePduRequest((short) 960);
            negotiatePduRequest.serialize(dos);
            channel.write(ByteBuffer.wrap(baos.toByteArray()));

            ackBuffer.clear();
            channel.read(ackBuffer);
            dis = new DataInputStream(new ByteArrayInputStream(ackBuffer.array()));
            message = parser.parse(dis);
            if(!(message instanceof S7ResponseMessage)) {
                throw new PlcProtocolException("Expecting S7 Response");
            }
            S7ResponseMessage negotiatePduResponse = (S7ResponseMessage) message;
            SetupCommunicationParameter setupCommunicationParameter =
                negotiatePduResponse.getS7Parameter(SetupCommunicationParameter.class);
            System.out.println("Negotiated PDU size of " + setupCommunicationParameter.getPduLength());
            System.out.println("Ready to communicate");
        } catch (UnknownHostException e) {
            throw new PlcIoException("Unknown Host " + hostName, e);
        } catch (ClosedChannelException e) {
            throw new PlcIoException("Channel already closed", e);
        } catch (IOException e) {
            throw new PlcIoException("IO Exception", e);
        }
    }

    public Message sendMessage(Message input) throws PlcException {
        try {
            ByteArrayOutputStream baos = new ByteArrayOutputStream(pduSize);
            DataOutputStream dos = new DataOutputStream(baos);
            input.serialize(dos);
            PcapngUtils.dumpPacket(ByteBuffer.wrap(baos.toByteArray()), baos.size(), "read-request.pcapng");
            channel.write(ByteBuffer.wrap(baos.toByteArray()));
            ByteBuffer ackBuffer = ByteBuffer.allocate(pduSize);
            int numBytes = channel.read(ackBuffer);
            PcapngUtils.dumpPacket(ackBuffer, numBytes, "read-response.pcapng");
            return parseMessage(ackBuffer.array());
        } catch (IOException e) {
            throw new PlcIoException("Error sending message", e);
        }
    }

    private Message parseMessage(byte[] buffer) throws PlcException {
        DataInputStream dis = new DataInputStream(new ByteArrayInputStream(buffer));
        return new S7Parser().parse(dis);
    }

}
