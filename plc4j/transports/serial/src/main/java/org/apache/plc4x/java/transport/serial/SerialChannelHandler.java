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
package org.apache.plc4x.java.transport.serial;

import com.fazecast.jSerialComm.SerialPort;
import com.fazecast.jSerialComm.SerialPortDataListener;
import com.fazecast.jSerialComm.SerialPortEvent;
import io.netty.buffer.ByteBuf;

import java.net.SocketAddress;
import java.nio.channels.SelectionKey;

/**
 * This is a wrapper mostly for testing {@link SerialChannel}, @{@link SerialPollingSelector},
 * {@link SerialSelectionKey}, @{@link SerialSelectorProvider} and @{@link SerialSocketChannel}.
 */
public abstract class SerialChannelHandler {

    private final SocketAddress address;

    public SerialChannelHandler(SocketAddress address) {
        this.address = address;
    }

    abstract boolean open();

    abstract String getIdentifier();

    /**
     * This method registers the Callback to the SelectionKey / {@link java.nio.channels.Selector}
     * which is necessary to notify the {@link java.nio.channels.Selector} about
     * available data.
     */
    abstract void registerSelectionKey(SerialSelectionKey selectionKey);

    public abstract void close();

    /**
     *
     * @param buf
     * @return Return the amoubnt of bytes written into the buffer
     */
    public abstract int read(ByteBuf buf);

    /**
     *
     * @return Number of bytes written to wire
     */
    public abstract int write(ByteBuf buf);

    public static class DummyHandler extends SerialChannelHandler {

        public static final DummyHandler INSTANCE = new DummyHandler(null);

        private SerialSelectionKey selectionKey;

        public DummyHandler(SocketAddress address) {
            super(address);
        }

        @Override
        public boolean open() {
            return true;
        }

        @Override
        public String getIdentifier() {
            return null;
        }

        @Override
        public void registerSelectionKey(SerialSelectionKey selectionKey) {
            this.selectionKey = selectionKey;
        }

        @Override
        public void close() {
            // NOOP
        }

        @Override
        public int read(ByteBuf buf) {
            buf.writeByte(1);
            return 1;
        }

        @Override
        public int write(ByteBuf buf) {
            System.out.println("Haha i wrote something");
            return 1;
        }

        public void fireEvent(int readyOp) {
            ((SerialPollingSelector) this.selectionKey.selector())
                .addEvent(new SerialPollingSelector.SelectorEvent(this.selectionKey, readyOp));
        }
    }


    public static class SerialPortHandler extends SerialChannelHandler {

        private SerialChannelConfig config;
        private SerialPort comPort;

        public SerialPortHandler(SocketAddress address, SerialChannelConfig config) {
            super(address);
            this.config = config;
            // Get the serial port described by the path/name in the address.
            comPort = SerialPort.getCommPort(((SerialSocketAddress) address).getIdentifier());
        }

        @Override
        public boolean open() {
            if (comPort.openPort()) {
                comPort.setComPortParameters(config.getBaudRate(), config.getDataBits(),
                    config.getStopBits(), config.getParityBits());
                return true;
            }
            return false;
        }

        @Override
        public String getIdentifier() {
            return comPort.getDescriptivePortName();
        }

        @Override
        public void registerSelectionKey(SerialSelectionKey selectionKey) {
            comPort.addDataListener(new SerialPortDataListener() {
                @Override
                public int getListeningEvents() {
                    return SerialPort.LISTENING_EVENT_DATA_AVAILABLE;
                }

                @Override
                public void serialEvent(SerialPortEvent event) {
                    ((SerialPollingSelector) selectionKey.selector())
                        .addEvent(new SerialPollingSelector.SelectorEvent(selectionKey, SelectionKey.OP_READ));
                }
            });
        }

        @Override
        public void close() {
            this.comPort.closePort();
        }

        @Override
        public int read(ByteBuf buf) {
            int bytesToRead = comPort.bytesAvailable();
            assert bytesToRead > 0;
            byte[] buffer = new byte[bytesToRead];
            comPort.readBytes(buffer, bytesToRead);
            buf.writeBytes(buffer);
            return bytesToRead;
        }

        @Override
        public int write(ByteBuf buf) {
            int expectedToWrite = buf.readableBytes();
            byte[] bytes = new byte[expectedToWrite];
            buf.readBytes(bytes);
            int bytesWritten = comPort.writeBytes(bytes, expectedToWrite);
            return bytesWritten;
        }
    }

}
