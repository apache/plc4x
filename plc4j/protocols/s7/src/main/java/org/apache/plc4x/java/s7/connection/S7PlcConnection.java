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

import org.apache.mina.core.buffer.IoBuffer;
import org.apache.mina.core.future.CloseFuture;
import org.apache.mina.core.future.ConnectFuture;
import org.apache.mina.core.session.IoSession;
import org.apache.mina.transport.socket.nio.NioSocketConnector;
import org.apache.plc4x.java.api.connection.AbstractPlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcIoException;
import org.apache.plc4x.java.isoontcp.mina.IsoOnTcpFilterAdapter;
import org.apache.plc4x.java.isotp.mina.IsoTPFilterAdapter;
import org.apache.plc4x.java.isotp.mina.model.tpdus.DisconnectRequestTpdu;
import org.apache.plc4x.java.isotp.mina.model.types.DisconnectReason;
import org.apache.plc4x.java.api.messages.PlcSimpleReadResponse;
import org.apache.plc4x.java.mina.PlcRequestContainer;
import org.apache.plc4x.java.api.messages.Address;
import org.apache.plc4x.java.api.messages.PlcSimpleReadRequest;
import org.apache.plc4x.java.s7.mina.Plc4XS7FilterAdapter;
import org.apache.plc4x.java.s7.mina.S7FilterAdapter;
import org.apache.plc4x.java.s7.mina.S7Handler;
import org.apache.plc4x.java.s7.mina.model.types.MemoryArea;
import org.apache.plc4x.java.s7.model.S7Address;
import org.apache.plc4x.java.s7.model.S7BitAddress;
import org.apache.plc4x.java.s7.model.S7DataBlockAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.UnknownHostException;
import java.util.Collections;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class S7PlcConnection extends AbstractPlcConnection implements PlcReader {

    private static final int ISO_ON_TCP_PORT = 102;

    private static final Pattern S7_DATABLOCK_ADDRESS_PATTERN = Pattern.compile("^DATABLOCK/(\\d{1,4})/(\\d{1,4})");
    private static final Pattern S7_ADDRESS_PATTERN = Pattern.compile("^(.*?)/(\\d{1,4})(?:/(\\d))?");

    private final static Logger logger = LoggerFactory.getLogger(S7PlcConnection.class);

    private final String hostName;
    private final int rack;
    private final int slot;

    private int pduSize;

    private IoSession session;

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
            InetAddress serverIPAddress = InetAddress.getByName(hostName);

            // Build the protocol stack for communicating with the s7 protocol.
            NioSocketConnector connector = new NioSocketConnector();
            connector.getFilterChain().addLast("iso-on-tcp", new IsoOnTcpFilterAdapter());
            connector.getFilterChain().addLast("iso-tp", new IsoTPFilterAdapter());
            connector.getFilterChain().addLast("s7", new S7FilterAdapter());
            connector.getFilterChain().addLast("plc4x-s7", new Plc4XS7FilterAdapter());

            // Create a future to make it possible to signal back as soon as the session
            // is completely initialized.
            CompletableFuture<Void> sessionSetupCompleteFuture = new CompletableFuture<>();
            connector.setHandler(new S7Handler() {
                @Override
                public void sessionOpened(IoSession session) throws Exception {
                    sessionSetupCompleteFuture.complete(null);
                }
            });

            // Connect to the PLC and wait till the session is created.
            // Pass in the attributes needed to parametrize the connection.
            ConnectFuture future = connector.connect(
                new InetSocketAddress(serverIPAddress, ISO_ON_TCP_PORT),
                (session, initFuture) -> {
                    session.setAttribute(IsoTPFilterAdapter.RACK_NO, (byte) rack);
                    session.setAttribute(IsoTPFilterAdapter.SLOT_NO, (byte) slot);
                    session.setAttribute(S7FilterAdapter.REQUESTED_PDU_SIZE, (short) pduSize);
                });

            // Wait until the TCP connection is established
            // (The negotiation on the higher protocols have not been handled then however)
            future.awaitUninterruptibly();
            session = future.getSession();

            // Wait till the "complete" method is called in the IoHandler
            sessionSetupCompleteFuture.get();

            logger.info("Session created");
        } catch (UnknownHostException e) {
            throw new PlcIoException("Unknown Host " + hostName, e);
        } catch (InterruptedException | ExecutionException e) {
            throw new PlcIoException(e);
        }
    }

    @Override
    public void close() throws Exception {
        // Send the PLC a message that the connection is being closed.
        DisconnectRequestTpdu disconnectRequest = new DisconnectRequestTpdu(
            (short) 0x0000, (short) 0x000F, DisconnectReason.NORMAL, Collections.emptyList(),
            IoBuffer.allocate(0).setAutoExpand(true));
        session.write(disconnectRequest);

        // Close the session itself.
        CloseFuture closeFuture = session.closeOnFlush();
        closeFuture.awaitUninterruptibly();
    }

    @Override
    public Address parseAddress(String addressString) throws PlcException {
        Matcher datablockAddressMatcher = S7_DATABLOCK_ADDRESS_PATTERN.matcher(addressString);
        if (datablockAddressMatcher.matches()) {
            int datablockNumber = Integer.valueOf(datablockAddressMatcher.group(1));
            int datablockByteOffset = Integer.valueOf(datablockAddressMatcher.group(2));
            return new S7DataBlockAddress((short) datablockNumber, (short) datablockByteOffset);
        }
        Matcher addressMatcher = S7_ADDRESS_PATTERN.matcher(addressString);
        if (!addressMatcher.matches()) {
            throw new PlcConnectionException(
                "Address string doesn't match the format '{area}/{offset}[/{bit-offset}]'");
        }
        MemoryArea memoryArea = MemoryArea.valueOf(addressMatcher.group(1));
        int byteOffset = Integer.valueOf(addressMatcher.group(2));
        if (addressMatcher.groupCount() == 4) {
            int bitOffset = Integer.valueOf(addressMatcher.group(3));
            return new S7BitAddress(memoryArea, (short) byteOffset, (byte) bitOffset);
        }
        return new S7Address(memoryArea, (short) byteOffset);
    }

    @Override
    public CompletableFuture<PlcSimpleReadResponse> read(PlcSimpleReadRequest readRequest) {
        CompletableFuture<PlcSimpleReadResponse> readFuture = new CompletableFuture<>();
        PlcRequestContainer<PlcSimpleReadRequest, PlcSimpleReadResponse> container =
            new PlcRequestContainer<>(readRequest, readFuture);
        session.write(container);
        return readFuture;
    }

}
