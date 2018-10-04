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
package org.apache.plc4x.java.ads.connection;

import io.netty.buffer.ByteBuf;
import io.netty.channel.embedded.EmbeddedChannel;
import io.netty.channel.jsc.JSerialCommDeviceAddress;
import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.commons.lang3.reflect.MethodUtils;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.serial.AmsSerialAcknowledgeFrame;
import org.apache.plc4x.java.ads.api.serial.AmsSerialFrame;
import org.apache.plc4x.java.ads.api.serial.types.*;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.base.connection.AbstractPlcConnection;
import org.apache.plc4x.java.base.connection.SerialChannelFactory;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.lang.reflect.InvocationTargetException;
import java.util.Arrays;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.doReturn;
import static org.mockito.Mockito.spy;

public class AdsSerialPlcConnectionTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(AdsSerialPlcConnectionTest.class);

    private AdsSerialPlcConnection SUT;

    @Before
    public void setUp() {
        SUT = AdsSerialPlcConnection.of("/dev/tty0", AmsNetId.of("0.0.0.0.0.0"), AmsPort.of(13));
    }

    @After
    public void tearDown() {
        SUT = null;
    }

    @Test
    public void initialState() {
        assertEquals(SUT.getTargetAmsNetId().toString(), "0.0.0.0.0.0");
        assertEquals(SUT.getTargetAmsPort().toString(), "13");
    }


    @Test
    public void testRead() throws Exception {
        prepareSerialSimulator();
        PlcReadRequest request = SUT.readRequestBuilder().get().addItem("test", "0/0:BYTE").build();
        CompletableFuture<PlcReadResponse> read = SUT.read(request);
        PlcReadResponse plcReadResponse = read.get(30, TimeUnit.SECONDS);
        assertNotNull(plcReadResponse);
    }

    private void prepareSerialSimulator() throws Exception {
        Field channelFactoryField = FieldUtils.getField(AbstractPlcConnection.class, "channelFactory", true);
        SerialChannelFactory serialChannelFactory = (SerialChannelFactory) channelFactoryField.get(SUT);
        SerialChannelFactory serialChannelFactorySpied = spy(serialChannelFactory);
        EmbeddedChannel embeddedChannel = new EmbeddedChannel(SUT.getChannelHandler(null));
        embeddedChannel.connect(new JSerialCommDeviceAddress("/dev/tty0"));
        doReturn(embeddedChannel).when(serialChannelFactorySpied).createChannel(any());
        channelFactoryField.set(SUT, serialChannelFactorySpied);
        SUT.connect();
        new SerialSimulator(embeddedChannel).start();
    }

    private class SerialSimulator extends Thread {

        private EmbeddedChannel embeddedChannel;

        private SimulatorState state = SimulatorState.RECEIVE_REQUEST;

        private byte[] currentInvokeId = new byte[0];

        public SerialSimulator(EmbeddedChannel embeddedChannel) {
            super("Serial Simulator");
            this.embeddedChannel = embeddedChannel;
        }

        @Override
        public void run() {
            while (true) {
                LOGGER.trace("in state {}. CurrentInvokeId: {}", state, currentInvokeId);
                switch (state) {
                    // Receiving state
                    case RECEIVE_REQUEST: {
                        LOGGER.info("Waiting for normal message");
                        ByteBuf outputBuffer;
                        while ((outputBuffer = embeddedChannel.readOutbound()) == null) {
                            LOGGER.trace("No buffer available yet");
                            if (!trySleep()) {
                                return;
                            }
                        }
                        int headerBytes = MagicCookie.NUM_BYTES + TransmitterAddress.NUM_BYTES + ReceiverAddress.NUM_BYTES + FragmentNumber.NUM_BYTES;
                        LOGGER.info("Skipping " + headerBytes + " bytes");
                        outputBuffer.skipBytes(headerBytes);
                        short dataLength = outputBuffer.readUnsignedByte();
                        LOGGER.info("Expect at least " + dataLength + "bytes");
                        while (outputBuffer.readableBytes() < dataLength) {
                            if (!trySleep()) {
                                return;
                            }
                        }
                        byte[] bytes = new byte[dataLength];
                        LOGGER.info("Read " + dataLength + "bytes. Having " + outputBuffer.readableBytes() + "bytes");
                        outputBuffer.readBytes(bytes);
                        currentInvokeId = Arrays.copyOfRange(bytes, 28, 32);
                        outputBuffer.skipBytes(CRC.NUM_BYTES);
                        LOGGER.info("Wrote Inbound");
                        state = SimulatorState.ACK_MESSAGE;
                        if (!trySleep()) {
                            return;
                        }
                    }
                    break;
                    case ACK_MESSAGE: {
                        ByteBuf byteBuf = AmsSerialAcknowledgeFrame.of(
                            TransmitterAddress.of((byte) 0x0),
                            ReceiverAddress.of((byte) 0x0),
                            FragmentNumber.of((byte) 0)
                        ).getByteBuf();
                        try {
                            MethodUtils.invokeMethod(byteBuf, true, "setRefCnt", 2);
                        } catch (NoSuchMethodException | IllegalAccessException | InvocationTargetException e) {
                            throw new RuntimeException(e);
                        }
                        embeddedChannel.writeOneInbound(byteBuf);
                        LOGGER.info("Acked Message");
                        state = SimulatorState.SEND_RESPONSE;
                    }
                    case SEND_RESPONSE: {
                        LOGGER.info("Sending data message");
                        ByteBuf byteBuf = AmsSerialFrame.of(
                            FragmentNumber.of((byte) 0),
                            UserData.of(
                                new byte[]{
                                    /*NetID Empfaenger */    (byte) 0xC0, (byte) 0xA8, 0x64, (byte) 0x9C, 0x01, 0x01,
                                    /*Portnummer       */    0x01, (byte) 0x80,
                                    /*NetID Sender     */    (byte) 0xC0, (byte) 0xA8, 0x64, (byte) 0xAE, 0x01, 0x01,
                                    /*Portnummer       */    0x21, 0x03,
                                    /*Response Lesen   */    0x02, 0x00,
                                    /*Status           */    0x05, 0x00,
                                    /*Anzahl Daten     */    0x0A, 0x00, 0x00, 0x00,
                                    /*Fehlercode       */    0x00, 0x00, 0x00, 0x00,
                                    /*InvokeID         */    currentInvokeId[0], currentInvokeId[1], currentInvokeId[2], currentInvokeId[3],
                                    /*Ergebnis         */    0x00, 0x00, 0x00, 0x00,
                                    /*Anzahl Daten     */    0x02, 0x00, 0x00, 0x00,
                                    /*Daten            */    (byte) 0xAF, 0x27
                                }
                            )
                        ).getByteBuf();
                        try {
                            MethodUtils.invokeMethod(byteBuf, true, "setRefCnt", 2);
                        } catch (NoSuchMethodException | IllegalAccessException | InvocationTargetException e) {
                            throw new RuntimeException(e);
                        }
                        embeddedChannel.writeOneInbound(byteBuf);
                        LOGGER.info("Wrote Inbound");
                        state = SimulatorState.WAIT_FOR_ACK;
                        if (!trySleep()) {
                            return;
                        }
                    }
                    break;
                    case WAIT_FOR_ACK: {
                        LOGGER.info("Waiting for ack message");
                        ByteBuf outputBuffer;
                        while ((outputBuffer = embeddedChannel.readOutbound()) == null) {
                            if (!trySleep()) {
                                return;
                            }
                        }
                        int headerBytes = MagicCookie.NUM_BYTES + TransmitterAddress.NUM_BYTES + ReceiverAddress.NUM_BYTES + FragmentNumber.NUM_BYTES;
                        LOGGER.info("Skipping " + headerBytes + " bytes");
                        outputBuffer.skipBytes(headerBytes);
                        short dataLength = outputBuffer.readUnsignedByte();
                        LOGGER.info("Expect " + dataLength + "bytes");
                        state = SimulatorState.DONE;
                        if (!trySleep()) {
                            return;
                        }
                    }
                    case DONE: {
                        LOGGER.info("Plc is Done. Goodbye");
                        return;
                    }
                    default:
                        throw new IllegalStateException("Illegal state number" + state);
                }
            }

        }

        private boolean trySleep() {
            try {
                TimeUnit.MILLISECONDS.sleep(10);
            } catch (InterruptedException e) {
                e.printStackTrace();
                Thread.currentThread().interrupt();
                return false;
            }
            return true;
        }
    }

    private enum SimulatorState {
        RECEIVE_REQUEST,
        ACK_MESSAGE,
        SEND_RESPONSE,
        WAIT_FOR_ACK,
        DONE
    }
}