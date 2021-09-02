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

import org.apache.commons.lang3.NotImplementedException;

import java.io.IOException;
import java.net.ProtocolFamily;
import java.nio.channels.*;
import java.nio.channels.spi.AbstractSelector;
import java.nio.channels.spi.SelectorProvider;

/**
 * TODO write comment
 *
 * @author julian
 * Created by julian on 2019-08-10
 */
class SerialSelectorProvider extends SelectorProvider {

    @Override
    public DatagramChannel openDatagramChannel() throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public DatagramChannel openDatagramChannel(ProtocolFamily family) throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public Pipe openPipe() throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public AbstractSelector openSelector() throws IOException {
        return new SerialPollingSelector(this);
    }

    @Override
    public ServerSocketChannel openServerSocketChannel() throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public SocketChannel openSocketChannel() throws IOException {
        return new SerialSocketChannel(this);
    }

}
