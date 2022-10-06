/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.transport.serial;

import io.netty.buffer.ByteBuf;

import java.net.SocketAddress;

public class DummyHandler extends SerialChannelHandler {

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
