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
package org.apache.plc4x.java.modbus.connection;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.embedded.EmbeddedChannel;
import io.netty.channel.jsc.JSerialCommDeviceAddress;
import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.commons.lang3.reflect.MethodUtils;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.base.connection.AbstractPlcConnection;
import org.apache.plc4x.java.base.connection.SerialChannelFactory;
import org.junit.After;
import org.junit.Before;
import org.junit.Ignore;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.lang.reflect.InvocationTargetException;
import java.util.Arrays;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

import static org.junit.Assert.assertNotNull;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.doReturn;
import static org.mockito.Mockito.spy;

@Ignore("Not yet implemented in modbus")
public class ModbusSerialPlcConnectionTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(ModbusSerialPlcConnectionTest.class);

    private ModbusSerialPlcConnection SUT;

    @Before
    public void setUp() {
        SUT = ModbusSerialPlcConnection.of("/dev/tty0", null);
    }

    @After
    public void tearDown() {
        SUT = null;
    }

    @Test
    public void testRead() throws Exception {
        prepareSerialSimulator();
        PlcReadRequest readRequest = SUT.readRequestBuilder().get().addItem("randomRead", "0/0").build();
        CompletableFuture<PlcReadResponse> read = SUT.read(readRequest);
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
                        // TODO
                        int headerBytes = 4711;
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
                        // TODO
                        outputBuffer.skipBytes(4711);
                        LOGGER.info("Wrote Inbound");
                        state = SimulatorState.ACK_MESSAGE;
                        if (!trySleep()) {
                            return;
                        }
                    }
                    break;
                    case ACK_MESSAGE: {
                        // TODO
                        ByteBuf byteBuf = Unpooled.buffer();
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
                        //TODO:
                        ByteBuf byteBuf = Unpooled.buffer();
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
                        //TODO:
                        int headerBytes = 4711;
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