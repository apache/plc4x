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

package org.apache.plc4x.java.transport.cotp;

import org.apache.plc4x.java.cotp.readwrite.*;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.transports.api.AsyncTransportInstance;
import org.apache.plc4x.java.spi.transports.api.BaseTransportInstance;
import org.apache.plc4x.java.spi.transports.api.RingBuffer;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.spi.utils.StaticHelper;
import org.apache.plc4x.java.transport.cotp.config.CotpTransportConfiguration;
import org.apache.plc4x.java.transport.tcp.TcpTransport;
import org.apache.plc4x.java.transport.tcp.TcpTransportInstance;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;
import java.util.function.Consumer;

/**
 * COTP Transport Instance - uses generated protocol classes for TPKT/COTP framing.
 * <p>
 * This implementation wraps TCP and uses the generated COTP protocol classes
 * to handle packet parsing and serialization, following the same pattern
 * that drivers use for protocol handling.
 * <p>
 * Implements AsyncTransportInstance to propagate data-available events from the
 * underlying TCP transport up to the driver layer, enabling event-driven I/O
 * without polling loops.
 */
public class CotpTransportInstance extends BaseTransportInstance<CotpTransportConfiguration> implements AsyncTransportInstance<CotpTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(CotpTransportInstance.class);

    private static final int TPKT_HEADER_SIZE = 4;
    private static final int DEFAULT_BUFFER_SIZE = 8192;

    private final TcpTransportInstance tcpTransport;
    private final RingBuffer payloadBuffer;
    private boolean connected = false;

    // Async support - propagate events from the TCP layer
    private volatile Runnable dataListener;
    private volatile Consumer<Throwable> disconnectListener;

    public CotpTransportInstance(String host, int port, CotpTransportConfiguration configuration, AuditLog auditLog)
            throws TransportException {
        super(configuration, auditLog);
        this.payloadBuffer = new RingBuffer(DEFAULT_BUFFER_SIZE);

        try {
            // Create underlying TCP transport
            TcpTransport tcpTransportFactory = new TcpTransport();
            this.tcpTransport = (TcpTransportInstance) tcpTransportFactory.createTransportInstance(
                host + ":" + port, configuration, auditLog);

            // Set up async event propagation from the TCP layer
            setupAsyncEventPropagation();

            // Perform COTP connection handshake
            performCotpConnection();

        } catch (Exception e) {
            String errorMsg = String.format("Failed to establish COTP connection to %s:%d - %s",
                host, port, e.getMessage());
            getAuditLog().write(AuditLogEventType.ERROR, "Error in constructor: " + errorMsg);
            throw new TransportException("Failed to establish COTP connection", e);
        }
    }

    /**
     * Sets up event propagation from the TCP transport to notify upper layers
     * when data becomes available. This enables event-driven I/O.
     * <p>
     * The TCP transport has already filled it's buffer before calling our callback,
     * so we just need to parse the available COTP packets and notify upward.
     */
    private void setupAsyncEventPropagation() {
        // Register a listener on TCP transport to parse COTP packets and notify
        // Note: TCP has already filled its ring buffer before calling this callback!
        tcpTransport.registerDataListener(() -> {
            try {
                // Parse any available COTP packets from TCP's already-filled buffer
                parseCotpPackets();

                // Notify our listener if we have parsed payload data
                if (payloadBuffer.availableForReading() > 0) {
                    Runnable listener = dataListener;
                    if (listener != null) {
                        listener.run();
                    }
                }
            } catch (TransportException e) {
                LOGGER.error("Error parsing COTP packets in async callback", e);
            }
        });

        LOGGER.debug("Set up async event propagation from TCP to COTP layer");
    }

    /**
     * Performs COTP connection handshake using generated protocol classes.
     */
    private void performCotpConnection() throws TransportException {
        LOGGER.info("Starting COTP connection handshake");

        try {
            // Build Connection Request using generated classes
            COTPPacketConnectionRequest connectionRequest = buildConnectionRequest();
            LOGGER.info("COTP CR: localTsap=0x{}, remoteTsap=0x{}",
                Integer.toHexString(getConfiguration().getLocalTsap()),
                Integer.toHexString(getConfiguration().getRemoteTsap()));

            TPKTPacket tpktRequest = new TPKTPacket(connectionRequest);

            // Serialize and send - use exact size to avoid extra bytes
            int packetSize = tpktRequest.getLengthInBytes();
            WriteBuffer writeBuffer = new WriteBufferByteBased(new byte[packetSize]);
            tpktRequest.serialize(writeBuffer);
            byte[] requestBytes = writeBuffer.getBytes();

            LOGGER.info("Sending COTP Connection Request ({} bytes): {}", requestBytes.length, StaticHelper.ENCODE_HEX(requestBytes));
            tcpTransport.write(requestBytes);

            // Wait for Connection Confirm
            long startTime = System.currentTimeMillis();
            while (System.currentTimeMillis() - startTime < getConfiguration().cotpConnectionTimeout) {
                int available = tcpTransport.getNumBytesAvailable();
                if (available >= TPKT_HEADER_SIZE + 1) {
                    // Peek at TPKT header to get full packet length
                    byte[] headerBytes = tcpTransport.peekReadableBytes(TPKT_HEADER_SIZE + 1);
                    int packetLength = TPKT_HEADER_SIZE + 1 + headerBytes[4];

                    if (available >= packetLength) {
                        // Read a complete packet
                        byte[] packetBytes = tcpTransport.read(packetLength);

                        // Parse using generated code
                        ReadBuffer readBuffer = new ReadBufferByteBased(packetBytes);
                        TPKTPacket tpktPacket = TPKTPacket.staticParse(readBuffer);
                        COTPPacket cotpPacket = tpktPacket.getPayload();

                        // Check for Connection Confirm
                        if (cotpPacket instanceof COTPPacketConnectionResponse) {
                            LOGGER.info("COTP connection established");

                            // Update the Tpdu size to the negotiated value.
                            cotpPacket.getParameters().stream()
                                .filter(COTPParameterTpduSize.class::isInstance)
                                .map(COTPParameterTpduSize.class::cast)
                                .findFirst()
                                .ifPresent(p ->
                                    getConfiguration().cotpTpduSize = p.getTpduSize().getSizeInBytes()
                                );

                            connected = true;
                            getAuditLog().write(AuditLogEventType.CONNECT, String.format(
                                "COTP connection established (localTsap=0x%04X, remoteTsap=0x%04X, tpduSize=%d)",
                                getConfiguration().getLocalTsap(), getConfiguration().getRemoteTsap(), getConfiguration().cotpTpduSize));
                            return;
                        } else if (cotpPacket instanceof COTPPacketTpduError) {
                            throw new TransportException("COTP connection rejected by remote");
                        }
                    }
                }
            }

            String errorMsg = String.format(
                "COTP connection timeout - no response from remote after %dms. " +
                "TSAP: 0x%04X → 0x%04X, TPDU size: %d bytes",
                getConfiguration().cotpConnectionTimeout,
                getConfiguration().getLocalTsap(),
                getConfiguration().getRemoteTsap(),
                getConfiguration().cotpTpduSize
            );
            throw new TransportException(errorMsg);

        } catch (Exception e) {
            String errorMsg = String.format(
                "COTP connection failed (TSAP: 0x%04X → 0x%04X, TPDU: %d bytes). " +
                "Check if PLC is reachable and configured to accept connections. " +
                "Error: %s",
                getConfiguration().getLocalTsap(),
                getConfiguration().getRemoteTsap(),
                getConfiguration().cotpTpduSize,
                e.getMessage()
            );
            throw new TransportException(errorMsg, e);
        }
    }

    /**
     * Builds a COTP Connection Request using generated protocol classes.
     */
    private COTPPacketConnectionRequest buildConnectionRequest() {
        // Build parameters
        List<COTPParameter> parameters = new ArrayList<>();

        // Source TSAP parameter
        int localTsap = getConfiguration().getLocalTsap();
        byte[] srcTsapBytes = new byte[]{
            (byte) ((localTsap >> 8) & 0xFF),
            (byte) (localTsap & 0xFF)
        };
        parameters.add(new COTPParameterCallingTsap(srcTsapBytes));

        // Destination TSAP parameter
        int remoteTsap = getConfiguration().getRemoteTsap();
        byte[] dstTsapBytes = new byte[]{
            (byte) ((remoteTsap >> 8) & 0xFF),
            (byte) (remoteTsap & 0xFF)
        };
        parameters.add(new COTPParameterCalledTsap(dstTsapBytes));

        // TPDU Size parameter
        COTPTpduSize tpduSize = COTPTpduSize.firstEnumForFieldSizeInBytes(getConfiguration().cotpTpduSize);
        parameters.add(new COTPParameterTpduSize(tpduSize));

        // Protocol class
        COTPProtocolClass protocolClass = COTPProtocolClass.enumForValue(
            (short) (getConfiguration().protocolClass << 4)
        );

        // Build Connection Request packet
        return new COTPPacketConnectionRequest(
            parameters,
            new byte[0],  // No payload for CR
            0x0000,       // Destination reference (0 for CR)
            0x000F,       // Source reference
            protocolClass
        );
    }

    @Override
    public boolean isOpen() {
        return connected && tcpTransport.isOpen();
    }

    /**
     * Parses any available COTP packets from the TCP buffer and extracts their payloads
     * into the COTP payload buffer. This assumes the TCP buffer is already filled.
     * <p>
     * This method is called both:
     * - From async callbacks (TCP already filled by selector)
     * - From fillBuffer() (after explicit TCP fillBuffer call)
     */
    private void parseCotpPackets() throws TransportException {
        if (!isOpen()) {
            return;
        }

        // Parse COTP packets and extract the payload into our buffer
        while (payloadBuffer.remainingForWriting() > 0) {
            int available = tcpTransport.getNumBytesAvailable();
            if (available < TPKT_HEADER_SIZE) {
                // Not enough data for TPKT header
                break;
            }

            try {
                // Peek at TPKT-header to get packet length
                byte[] headerBytes = tcpTransport.peekReadableBytes(TPKT_HEADER_SIZE);
                int packetLength = ((headerBytes[2] & 0xFF) << 8) | (headerBytes[3] & 0xFF);

                if (available < packetLength) {
                    // Not enough data for a complete packet
                    break;
                }

                // Read and consume the complete TPKT packet from TCP
                byte[] packetBytes = tcpTransport.read(packetLength);

                // Parse using generated code
                ReadBuffer readBuffer = new ReadBufferByteBased(packetBytes);
                TPKTPacket tpktPacket = TPKTPacket.staticParse(readBuffer);
                COTPPacket cotpPacket = tpktPacket.getPayload();

                // Extract the payload and write to our buffer
                byte[] payload = cotpPacket.getPayload();
                if (payload != null && payload.length > 0) {
                    // Check if we have space in the buffer
                    if (payloadBuffer.remainingForWriting() >= payload.length) {
                        payloadBuffer.write(payload);
                        LOGGER.trace("Extracted {} bytes of payload from COTP packet ({} bytes total)",
                            payload.length, packetLength);
                    } else {
                        // Buffer full, can't process more packets right now
                        LOGGER.trace("Payload buffer full, {} bytes remaining",
                            payloadBuffer.remainingForWriting());
                        break;
                    }
                }

            } catch (Exception e) {
                LOGGER.error("Error parsing COTP packet", e);
                throw new TransportException("Failed to parse COTP packet", e);
            }
        }
    }

    @Override
    public int getNumBytesAvailable() {
        if (!isOpen()) {
            return 0;
        }

        // Return available payload bytes from our buffer
        return payloadBuffer.availableForReading();
    }

    @Override
    public byte[] peekReadableBytes(int numBytes) throws TransportException {
        if (!isOpen()) {
            throw new TransportException("Transport is not open");
        }

        if (numBytes <= 0) {
            return new byte[0];
        }

        // Fill buffer to ensure we have data
        int available = getNumBytesAvailable();
        if (available < numBytes) {
            throw new TransportException(
                String.format("Requested %d bytes but only %d available", numBytes, available));
        }

        // Peek without consuming
        return payloadBuffer.peek(numBytes);
    }

    @Override
    public byte[] read(int numBytes) throws TransportException {
        if (!isOpen()) {
            throw new TransportException("Transport is not open");
        }

        if (numBytes <= 0) {
            return new byte[0];
        }

        // Fill buffer to ensure we have data
        int available = getNumBytesAvailable();
        if (available < numBytes) {
            throw new TransportException(
                String.format("Requested %d bytes but only %d available", numBytes, available));
        }

        // Read and consume bytes from the buffer
        byte[] result = payloadBuffer.read(numBytes);
        LOGGER.trace("Read {} bytes from COTP payload buffer", numBytes);
        return result;
    }

    @Override
    public void write(byte[] bytes) throws TransportException {
        if (!isOpen()) {
            throw new TransportException("Transport is not open");
        }

        if (bytes == null || bytes.length == 0) {
            return;
        }

        try {
            // Build COTP Data packet using generated classes
            COTPPacketData dataPacket = new COTPPacketData(
                new ArrayList<>(),  // No parameters for data packets
                bytes,              // Payload
                true,               // EOT (End of Transmission)
                (byte) 0            // TPDU reference number
            );

            // Wrap in TPKT
            TPKTPacket tpktPacket = new TPKTPacket(dataPacket);
            int lengthInBytes = tpktPacket.getLengthInBytes();

            // Serialize using generated code
            WriteBuffer writeBuffer = new WriteBufferByteBased(new byte[lengthInBytes]);
            tpktPacket.serialize(writeBuffer);
            byte[] packetBytes = writeBuffer.getBytes();

            // Send data via TCP
            tcpTransport.write(packetBytes);

            LOGGER.trace("Sent {} bytes of payload ({} bytes total with TPKT/COTP framing)",
                bytes.length, packetBytes.length);

        } catch (Exception e) {
            throw new TransportException("Failed to write bytes", e);
        }
    }

    @Override
    public void close() throws TransportException {
        if (!connected) {
            return;
        }

        try {
            // Send Disconnect Request using generated classes
            LOGGER.debug("Sending COTP Disconnect Request");

            COTPPacketDisconnectRequest disconnectRequest = new COTPPacketDisconnectRequest(
                new ArrayList<>(),
                new byte[0],
                0x0000,  // Destination reference
                0x0001,  // Source reference
                COTPProtocolClass.CLASS_0
            );

            TPKTPacket tpktPacket = new TPKTPacket(disconnectRequest);
            WriteBuffer writeBuffer = new WriteBufferByteBased(new byte[32]);
            tpktPacket.serialize(writeBuffer);

            tcpTransport.write(writeBuffer.getBytes());

            connected = false;

        } catch (Exception e) {
            LOGGER.warn("Error sending disconnect request", e);
        } finally {
            tcpTransport.close();
            LOGGER.debug("COTP transport closed");
            getAuditLog().write(AuditLogEventType.CLOSE, "COTP transport closed");
        }
    }

    // ========== AsyncTransportInstance Implementation ==========

    @Override
    public void registerDataListener(Runnable listener) {
        this.dataListener = listener;
        LOGGER.debug("Data listener registered on COTP transport");
    }

    @Override
    public void removeDataListener() {
        this.dataListener = null;
        LOGGER.debug("Data listener removed from COTP transport");
    }

    @Override
    public void registerDisconnectListener(Consumer<Throwable> listener) {
        this.disconnectListener = listener;
        // Propagate to TCP layer
        tcpTransport.registerDisconnectListener(this::onTcpDisconnected);
        LOGGER.debug("Disconnect listener registered on COTP transport");
    }

    @Override
    public void removeDisconnectListener() {
        this.disconnectListener = null;
        tcpTransport.removeDisconnectListener();
        LOGGER.debug("Disconnect listener removed from COTP transport");
    }

    /**
     * Called when the underlying TCP transport is disconnected.
     * Propagates the disconnect event to the COTP layer's listener.
     */
    private void onTcpDisconnected(Throwable cause) {
        connected = false;
        Consumer<Throwable> listener = disconnectListener;
        if (listener != null) {
            try {
                listener.accept(cause);
            } catch (Exception e) {
                LOGGER.error("Error in disconnect listener", e);
            }
        }
    }

}
