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

class IdleSerialChannelHandler extends SerialChannelHandler {

    public IdleSerialChannelHandler() {
        super(null);
    }

    @Override
    public boolean open() {
        return true;
    }

    @Override
    String getIdentifier() {
        return "Emulated Com Port";
    }

    @Override
    void registerSelectionKey(SerialSelectionKey selectionKey) {
        /*
         * We can ignore registration as we will neven trigger events
         */
    }

    @Override
    public void close() {
    }

    @Override
    public int read(ByteBuf buf) {
        return 0;
    }

    @Override
    public int write(ByteBuf buf) {
        throw new UnsupportedOperationException();
    }
}
